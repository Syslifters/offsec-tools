import json
import csv
from typing import List
from shareranger.models.host import Host
from shareranger.models.match import Match
from shareranger.logging.logger import log_success, log_status, log_error


def save_scan_results(hosts: List[Host], filepath: str):
    try:
        with open(filepath, "w", encoding="utf-8") as f:
            json.dump([h.to_dict() for h in hosts], f, indent=2)
            log_success(f"Saved scan results to: {filepath}")
    except Exception as e:
        log_error(f"Failed to save scan results to: {e}")


def load_scan_results(filepath: str) -> List[Host]:
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            data = json.load(f)
            result_hosts = [Host.from_dict(entry) for entry in data]
            log_status(f"Loaded {len(result_hosts)} hosts from {filepath}")
    except Exception as e:
        log_error(f"Failed to load scan results from: {e}")
    return result_hosts


def export_matches_to_csv(hosts: List[Host], filepath: str):
    try:
        with open(filepath, mode="w", newline="", encoding="utf-8") as f:
            writer = csv.DictWriter(
                f,
                fieldnames=[
                    "action",
                    "severity",
                    "type",
                    "host",
                    "share",
                    "path",
                    "id",
                    "location",
                    "pattern",
                    "snippet",
                ],
            )
            writer.writeheader()

            for host in hosts:
                for share in host.shares:
                    for match in share.matches:
                        match_dict = match.to_dict()
                        writer.writerow(
                            {
                                "action": match_dict.get("action", ""),
                                "severity": match_dict.get("severity", ""),
                                "type": match_dict.get("type", ""),
                                "host": host.name,
                                "share": share.name,
                                "path": match_dict.get("unc_path", ""),
                                "id": match_dict.get("id", ""),
                                "location": match_dict.get("location", ""),
                                "pattern": match_dict.get("pattern", ""),
                                "snippet": match_dict.get("snippet", ""),
                            }
                        )

        log_success(f"Exported findings to CSV: {filepath}")
    except Exception as e:
        log_error(f"Failed to export matches to CSV: {e}")
