import os
import yaml
from shareranger.config import config
from shareranger.models.rule import Rule, RuleAction, RuleType
from shareranger.logging.logger import log_status, log_debug, log_warning


def load_rules_from_dir(rule_dir: str) -> list[Rule]:
    rules = []
    seen_ids = set()

    for root, _, files in os.walk(rule_dir):
        for file in files:
            if file.endswith(".yaml") or file.endswith(".yml"):
                path = os.path.join(root, file)
                try:
                    with open(path, "r", encoding="utf-8") as f:
                        loaded = yaml.safe_load(f)
                        if not isinstance(loaded, list):
                            log_warning(
                                f"Invalid YAML structure in {file}: expected a list of rules"
                            )
                            continue
                        for raw in loaded:
                            rule = Rule.from_dict(raw, file)
                            if rule is None:
                                continue
                            if rule.id in seen_ids:
                                log_warning(
                                    f"Skipping duplicate rule ID '{rule.id}' in {file}"
                                )
                                continue
                            seen_ids.add(rule.id)
                            rules.append(rule)
                except Exception as e:
                    log_warning(f"Failed to load rules from {file}: {e}")

    return rules


def sort_rules(rules: list[Rule]) -> list[Rule]:
    return sorted(
        rules,
        key=lambda r: (r.action == RuleAction.MATCH, r.action == RuleAction.KEEP, r.id),
    )


def build_relay_tree(
    rule: Rule, rule_map: dict[str, Rule], used_in_tree: set[str], visited=None
):
    if visited is None:
        visited = set()

    if rule.id in visited:
        return None
    visited.add(rule.id)

    if rule.relay:
        children = []
        for rid in rule.relay:
            used_in_tree.add(rid)
            target = rule_map.get(rid)
            if target:
                child_tree = build_relay_tree(target, rule_map, used_in_tree, visited)
                if child_tree:
                    children.append(child_tree)
        rule.children = sort_rules(children)
    return rule


def print_rule_tree(rule: Rule, prefix=""):
    line = f"{prefix}{rule.id} ({rule.action.value})"
    log_debug(line)
    if rule.children:
        for i, child in enumerate(rule.children):
            connector = "└── " if i == len(rule.children) - 1 else "├── "
            print_rule_tree(child, prefix + connector)


def load_and_resolve_rules(rule_dir: str):
    rules = load_rules_from_dir(rule_dir)
    rule_map = {r.id: r for r in rules}
    used_in_tree = set()

    # Categorize rules by type
    categorized = {RuleType.SHARE: [], RuleType.DIRECTORY: [], RuleType.FILE: []}
    for rule in rules:
        categorized[rule.type].append(rule)

    resolved = {}
    for rtype in [RuleType.SHARE, RuleType.DIRECTORY, RuleType.FILE]:
        all_rules = categorized[rtype]
        rule_map = {r.id: r for r in all_rules}
        used_in_tree.clear()
        relay_trees = []

        for rule in all_rules:
            if rule.relay:
                tree = build_relay_tree(rule, rule_map, used_in_tree)
                if tree:
                    relay_trees.append(tree)

        relay_trees = sort_rules(relay_trees)
        flat_rules = [
            r
            for r in all_rules
            if r.relay is None
            and r.action != RuleAction.MATCH
            and r.id not in used_in_tree
        ]
        flat_rules = sort_rules(flat_rules)

        resolved[rtype.value] = {"flat": flat_rules, "trees": relay_trees}

        log_status(f"Loaded {len(all_rules)} {rtype.value} rules")

        for rule in flat_rules:
            log_debug(f"{rule.id} ({rule.action.value})")
        for relay in relay_trees:
            print_rule_tree(relay)

    return resolved
