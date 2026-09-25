import sys
import re
import xml.etree.ElementTree as ET
from typing import List, Dict
import dns.resolver
from ldap3 import Server, Connection, SASL, KERBEROS, BASE, SUBTREE
import win32com.client
import win32security
from shareranger.models.host import Host
from shareranger.logging.logger import log_success, log_status, log_warning, log_error


def get_base_dn():
    try:
        ad_system = win32com.client.Dispatch("ADSystemInfo")
        domain = ad_system.DomainDNSName
        parts = domain.split(".")
        return ",".join(f"DC={part}" for part in parts)
    except Exception:
        return None


def get_domain_from_base_dn(base_dn: str) -> str:
    parts = [part[3:] for part in base_dn.split(",") if part.startswith("DC=")]
    return ".".join(parts)


def get_default_dc(domain_dns_name: str):
    try:
        srv_query = f"_ldap._tcp.dc._msdcs.{domain_dns_name}"
        answers = dns.resolver.resolve(srv_query, "SRV")
        dcs = sorted(
            [str(rdata.target).rstrip(".") for rdata in answers],
            key=lambda x: x.lower(),
        )
        return dcs[0] if dcs else None
    except Exception as e:
        log_warning(f"Could not resolve DC via DNS SRV: {e}")
        return None


def establish_ldap_connection(dc: str, use_ssl: bool = True):
    port = 636 if use_ssl else 389
    server = Server(dc, port=port, use_ssl=use_ssl, get_info=BASE)
    conn = Connection(server, authentication=SASL, sasl_mechanism=KERBEROS)

    try:
        if conn.bind():
            if use_ssl:
                log_success(f"Secure LDAPS connection established to {dc}")
            else:
                log_success(f"Plain LDAP connection established to {dc}")
            return conn
        elif use_ssl:
            raise RuntimeError("LDAPS bind failed")
        else:
            raise RuntimeError("LDAP bind failed")
    except Exception as e:
        if use_ssl:
            log_warning(f"LDAPS bind failed: {e}. Falling back to plain LDAP.")
            return establish_ldap_connection(dc, use_ssl=False)
        else:
            log_error(f"LDAP bind failed: {e}. Aborting.")
            return None


def get_ldap_connection(dc: str = None, base_dn: str = None):
    if not base_dn:
        base_dn = get_base_dn()
        if base_dn:
            log_status(f"No Base DN provided, using detected Base DN: {base_dn}")
        else:
            log_error(
                "Base DN could not be detected automatically. Please specify it explicitly."
            )
            sys.exit(1)

    if not dc:
        base_domain = ".".join(
            part[3:] for part in base_dn.split(",") if part.startswith("DC=")
        )
        dc = get_default_dc(base_domain)

    if dc:
        log_status(f"No DC provided, resolved DC via DNS: {dc}")
    else:
        log_error("Could not resolve a Domain Controller. Aborting.")
        sys.exit(1)

    conn = establish_ldap_connection(dc)
    if not conn:
        sys.exit(1)

    return conn, base_dn


def ad_discover_hosts(conn, base_dn: str):
    conn.search(
        search_base=base_dn,
        search_filter="(objectClass=computer)",
        attributes=["dNSHostName"],
    )
    return [Host(name=str(entry["dNSHostName"])) for entry in conn.entries]


def ad_discover_dfs_namespaces(conn, base_dn: str):
    conn.search(
        search_base=base_dn,
        search_filter="(objectClass=msDFS-NamespaceAnchor)",
        attributes=["cn"],
    )
    return [str(entry.cn) for entry in conn.entries]


def normalize_dfs_target(unc_path: str, base_dn: str) -> str:
    domain = get_domain_from_base_dn(base_dn)
    parts = unc_path.strip("\\").split("\\", 1)
    if len(parts) != 2:
        return unc_path  # Not a valid UNC path
    host, share = parts
    if "." not in host:
        host = f"{host}.{domain}"
    return f"\\\\{host}\\{share}".lower()


def ad_discover_dfs_links_and_targets(
    conn: Connection, base_dn: str, namespace: str
) -> List[Dict]:
    results = []
    search_base = f"CN={namespace},CN=Dfs-Configuration,CN=System,{base_dn}"

    try:
        conn.search(
            search_base=search_base,
            search_filter="(objectClass=msDFS-Linkv2)",
            search_scope=SUBTREE,
            attributes=["cn", "msDFS-LinkPathv2", "msDFS-TargetListv2"],
        )

        for entry in conn.entries:
            data = entry.entry_attributes_as_dict
            name = data.get("cn", ["unknown"])[0]
            path = data.get("msDFS-LinkPathv2", [""])[0]
            targets = []

            if "msDFS-TargetListv2" in data:
                try:
                    binary = data["msDFS-TargetListv2"]
                    if isinstance(binary, list):
                        binary = binary[0]

                    xml_string = binary.decode("utf-16-le", errors="ignore")

                    # Strip BOM if present
                    if xml_string.startswith("\ufeff"):
                        xml_string = xml_string.lstrip("\ufeff")

                    root = ET.fromstring(xml_string)
                    raw_targets = [
                        t.text for t in root.findall(".//{*}target") if t.text
                    ]
                    targets = [normalize_dfs_target(t, base_dn) for t in raw_targets]

                except Exception as e:
                    log_warning(f"Failed to parse DFSv2 XML for link '{name}': {e}")

            results.append({"name": name, "path": path, "targets": targets})

    except Exception as e:
        log_warning(
            f"Failed to resolve DFS links via LDAP for namespace {namespace}: {e}"
        )

    return results


# HACKY
def ldap_token_groups(ldap_conn, sid_str) -> set[str]:
    filt = f"(objectSid:1.2.840.113556.1.4.146:={sid_str})"
    ldap_conn.search(
        "DC=lab,DC=local", filt, attributes=["distinguishedName", "tokenGroups"]
    )

    entry = ldap_conn.entries[0]  # should be exactly one
    dn = entry.distinguishedName.value
    token = {
        win32security.ConvertSidToStringSid(raw) for raw in entry.tokenGroups.values
    }
    return {
        win32security.ConvertSidToStringSid(raw)
        for raw in ldap_conn.entries[0].tokenGroups.values
    }
