import os
from concurrent.futures import ThreadPoolExecutor, as_completed
from threading import Lock
from shareranger.rules.matcher import match_rules
from shareranger.models.rule import RuleAction, RuleType
from shareranger.logging.logger import log_status
from shareranger.config import config
from shareranger.logging import events

# Lock to safely modify share.matches across threads
matches_lock = Lock()


def process_file(share, file_path, file_rules):
    file_name = os.path.basename(file_path)
    file_ext = os.path.splitext(file_path)[1]

    action, rule_matches = match_rules(
        flat_rules=file_rules.get("flat", []),
        relay_trees=file_rules.get("trees", []),
        type=RuleType.FILE,
        file_name=file_name,
        unc_path=file_path,
        file_ext=file_ext,
        share_name=share.name,
    )

    if action == RuleAction.KEEP:
        keep_matches = [m for m in rule_matches if m.rule.action == RuleAction.KEEP]
        if keep_matches:
            with matches_lock:
                share.matches.extend(keep_matches)
                for k in keep_matches:
                    events.emit_match(share._host_name, share.name, k.to_dict())


def walk_share(share, directory_rules, file_rules):
    log_status(f"Starting share scan: {share.unc_path}")

    with ThreadPoolExecutor(max_workers=config.args.file_threads) as executor:
        futures = []

        for root, dirs, files in os.walk(share.unc_path):
            action, rule_matches = match_rules(
                flat_rules=directory_rules.get("flat", []),
                relay_trees=directory_rules.get("trees", []),
                type=RuleType.DIRECTORY,
                unc_path=root,
                share_name=share.name,
            )

            if action == RuleAction.DISCARD:
                dirs.clear()
                continue
            elif action == RuleAction.KEEP:
                keep_matches = [
                    m for m in rule_matches if m.rule.action == RuleAction.KEEP
                ]
                with matches_lock:
                    share.matches.extend(keep_matches)
                    for k in keep_matches:
                        events.emit_match(share._host_name, share.name, k.to_dict())
                dirs.clear()
                continue

            for file in files:
                file_path = os.path.join(root, file)
                futures.append(
                    executor.submit(process_file, share, file_path, file_rules)
                )

        for future in as_completed(futures):
            future.result()

    log_status(f"Finished share scan: {share.unc_path}")
