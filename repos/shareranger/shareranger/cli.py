import os
import signal
from concurrent.futures import ThreadPoolExecutor

from shareranger.config import config
from shareranger.models.share import Share
from shareranger.models.host import Host
from shareranger.models.rule import RuleType
from shareranger.rules.loader import load_and_resolve_rules
from shareranger.discovery.ad import (
    get_ldap_connection,
    get_domain_from_base_dn,
    ad_discover_hosts,
    ad_discover_dfs_namespaces,
    ad_discover_dfs_links_and_targets,
)
from shareranger.discovery.smb import (
    scan_shares,
    process_share,
    check_smb_online,
)
from shareranger.logging import events
from shareranger.walker.fswalker import walk_share
from shareranger.logging.logger import init_logging, log_status, log_debug


def _hard_exit(signum, frame):
    # 130 is the conventional exit code for SIGINT; 143 for SIGTERM.
    code = 130 if signum == signal.SIGINT else 143
    log_status(f"Received stop signal, exiting immediately...")
    os._exit(code)


# Exit immediately on Ctrl+C and kill signals
signal.signal(signal.SIGINT, _hard_exit)  # Ctrl+C
signal.signal(signal.SIGTERM, _hard_exit)  # kill
try:
    signal.signal(signal.SIGPIPE, signal.SIG_DFL)  # optional: avoid BrokenPipe noise
except AttributeError:
    pass  # not on Windows

# On Windows, also handle Ctrl+Break for good measure
if hasattr(signal, "SIGBREAK"):
    signal.signal(signal.SIGBREAK, _hard_exit)


def run_phase1_discovery():
    # If --local-path is provided, skip Phase 1 and run Phase 2 directly on that folder
    if config.args.local_path:
        share = Share(
            name=os.path.basename(config.args.local_path),
            unc_path=config.args.local_path,
        )
        host = Host(name="localhost")
        host.add_share(share)
        events.emit_host(host)
        events.emit_share(host.name, share)
        return [host]

    conn, base_dn = get_ldap_connection(config.args.dc, config.args.base_dn)
    domain = get_domain_from_base_dn(base_dn)

    # Discover DFS namespaces and links
    dfs_namespaces = ad_discover_dfs_namespaces(conn, base_dn)
    dfs_targets = {}
    for ns in dfs_namespaces:
        unc_path = f"\\\\{domain}\\{ns}"
        links = ad_discover_dfs_links_and_targets(conn, base_dn, ns)
        for link in links:
            for target in link["targets"]:
                dfs_targets[target.lower()] = unc_path

    log_status(
        f"Discovered {len(dfs_namespaces)} DFS namespace(s) and {len(dfs_targets)} unique DFS target(s) in Active Directory."
    )

    # Discover SMB hosts
    discovered_hosts = ad_discover_hosts(conn, base_dn)
    log_status(f"Discovered {len(discovered_hosts)} host(s) in Active Directory.")

    # Scan shares per host
    def process_host(host):
        check_smb_online(host)
        if not host.online:
            return host
        if host.name.lower() != domain.lower():
            events.emit_host(host)
            share_infos = scan_shares(host)
            for info in share_infos:
                share = Share(
                    name=info["netname"],
                    unc_path=str(f"\\\\{host.name}\\{info['netname']}"),
                )
                if share.unc_path.lower() in dfs_targets:
                    log_debug(
                        f"Skipping DFS target {share.unc_path} covered by DFS namespace {dfs_targets[share.unc_path.lower()]}."
                    )
                    continue
                host.add_share(share)
        return host

    with ThreadPoolExecutor(max_workers=config.args.host_threads) as executor:
        result_hosts = list(executor.map(process_host, discovered_hosts))

    return result_hosts


def run_phase2_fswalking(hosts, rules):
    total_hosts = len(hosts)
    total_shares = sum(len(h.shares) for h in hosts)
    log_status(f"Scanning {total_hosts} host(s) and {total_shares} share(s)...")

    def process_and_walk(host, share):
        processed_share = process_share(host, share, rules[RuleType.SHARE])
        events.emit_share(host.name, processed_share)
        if processed_share.accessible and not processed_share.matches:
            walk_share(processed_share, rules[RuleType.DIRECTORY], rules[RuleType.FILE])

    with ThreadPoolExecutor(max_workers=config.args.share_threads) as executor:
        for host in hosts:
            for share in host.shares:
                executor.submit(process_and_walk, host, share)


def main():
    config.parse_args()

    init_logging(config.args.verbose)
    rules = load_and_resolve_rules(config.args.rule_dir)

    events.start(
        config.DEFAULT_RESULT_FILE,
        config.DEFAULT_RESULT_FILE_CSV if config.args.export_csv else None,
    )

    try:
        hosts = run_phase1_discovery()
        run_phase2_fswalking(hosts, rules)
    finally:
        events.stop()


if __name__ == "__main__":
    main()
