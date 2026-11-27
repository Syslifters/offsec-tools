import re
import os
from typing import Optional
from shareranger.config import config
from shareranger.discovery.ad import get_ldap_connection, ldap_token_groups
from shareranger.models.ace import ACE
from shareranger.models.rule import Rule, RuleAction, RuleType, RuleLocation
from shareranger.models.match import Match
from shareranger.models.share import Share
from shareranger.logging.logger import log_match, log_debug


def match_rules(
    flat_rules: list[Rule],
    relay_trees: list[Rule],
    type: RuleType,
    unc_path=None,
    file_name=None,
    file_ext=None,
    share_name=None,
):
    """
    Evaluates rules against a file or directory and returns the final action
    and all matching Rule objects.
    """
    matches = []
    visited = set()
    early_discard = False

    def get_target(location: RuleLocation):
        if location == RuleLocation.FILE_PATH:
            return unc_path
        elif location == RuleLocation.FILE_NAME:
            return file_name
        elif location == RuleLocation.FILE_EXTENSION:
            return file_ext
        elif location == RuleLocation.SHARE_NAME:
            return share_name
        elif location == RuleLocation.FILE_CONTENT:
            try:
                if unc_path and os.path.getsize(unc_path) <= config.args.max_bytes:
                    with open(unc_path, "rb") as f:
                        return f.read(config.args.max_bytes).decode(
                            "utf-8", errors="ignore"
                        )
                else:
                    log_debug(
                        f"Skipping FileContent match due to file size: {unc_path}"
                    )
            except Exception:
                return None
        return None

    def extract_snippet(text, match):
        start = max(match.start() - config.args.snippet_length, 0)
        end = min(match.end() + config.args.snippet_length, len(text))
        return text[start:end].replace("\n", " ").replace("\r", " ").strip()

    def match_single_rule(rule: Rule, visited):
        nonlocal early_discard

        if rule.id in visited or early_discard:
            return
        visited.add(rule.id)

        target = get_target(rule.location)
        if not target:
            return

        verbose = rule.action != RuleAction.KEEP

        for pattern in rule.patterns:
            regex = re.compile(pattern, re.IGNORECASE)
            match = regex.search(target)
            if match:
                risky_aces = []
                if rule.ntfs:
                    for sid in ACE.RISKY_SIDS:
                        try:
                            ace = ACE(path=unc_path, sid=sid)
                            if rule.ntfs.get("readable") and not ace.has_read:
                                continue
                            if rule.ntfs.get("writable") and not ace.has_write:
                                continue
                            risky_aces.append(ace)
                        except Exception as e:
                            log_debug(f"Error processing ACE for SID {sid}: {e}")
                            continue
                    if risky_aces:
                        log_debug(
                            f"Matched {len(risky_aces)} risky ACE(s) in rule {rule.id}"
                        )
                    else:
                        break  # no risky ACEs found, skip this match

                snippet = extract_snippet(target, match)
                log_match(
                    rule.action.value,
                    type.value,
                    unc_path,
                    rule.id,
                    rule.severity.value,
                    rule.location.value,
                    pattern,
                    snippet,
                    verbose,
                )
                matches.append(
                    Match(
                        rule=rule,
                        unc_path=unc_path,
                        pattern=pattern,
                        snippet=snippet,
                        insecure_aces=risky_aces,
                    )
                )

                if rule.action == RuleAction.DISCARD:
                    early_discard = True
                    return
                for child in rule.children:
                    match_single_rule(child, visited)
                break

    for rule in flat_rules:
        if rule.location in (
            RuleLocation.FILE_PATH,
            RuleLocation.FILE_NAME,
            RuleLocation.FILE_EXTENSION,
            RuleLocation.SHARE_NAME,
            RuleLocation.FILE_CONTENT,
        ):
            match_single_rule(rule, visited)
            if early_discard:
                return RuleAction.DISCARD, matches

    for relay in relay_trees:
        match_single_rule(relay, visited)
        if early_discard:
            return RuleAction.DISCARD, matches

    has_keep = any(m.rule.action == RuleAction.KEEP for m in matches)
    only_relay = matches and all(m.rule.action == RuleAction.MATCH for m in matches)

    if has_keep:
        return RuleAction.KEEP, matches
    elif only_relay:
        return RuleAction.MATCH, matches
    else:
        return None, matches  # nothing matched
