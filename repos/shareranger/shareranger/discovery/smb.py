import win32net
import socket

from shareranger.models.rule import RuleAction, RuleType
from shareranger.rules.matcher import match_rules
from shareranger.logging.logger import (
    log_success,
    log_status,
    log_warning,
    log_error,
)
from shareranger.logging import events


def check_smb_online(host):
    try:
        socket.create_connection((host.name, 445), timeout=3)
        host.online = True
        log_success(f"Host online: {host.name}")
    except Exception as e:
        host.online = False
        log_warning(f"Host unreachable: {host.name} ({e})")


def scan_shares(host):
    try:
        shares, _, _ = win32net.NetShareEnum(host.name, 1)
        log_status(f"Discovered {len(shares)} shares on {host.name}")
        return shares
    except Exception as e:
        log_error(f"Failed to list shares on {host.name}: {e}")
        return []


def process_share(host, share, share_rules):
    if not share.accessible:
        return share

    action, rule_matches = match_rules(
        share_rules.get("flat", []),
        share_rules.get("trees", []),
        type=RuleType.SHARE,
        unc_path=share.unc_path,
        share_name=share.name,
    )
    if action == RuleAction.DISCARD:
        discards = [m for m in rule_matches if m.rule.action == RuleAction.DISCARD]
        share.matches.extend(discards)
        for d in discards:
            events.emit_match(host.name, share.name, d.to_dict())
    elif action == RuleAction.KEEP:
        keeps = [m for m in rule_matches if m.rule.action == RuleAction.KEEP]
        share.matches.extend(keeps)
        for k in keeps:
            events.emit_match(host.name, share.name, k.to_dict())

    return share
