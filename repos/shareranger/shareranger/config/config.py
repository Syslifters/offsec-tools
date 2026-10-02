import argparse
from datetime import datetime
import sys
import os

# Global settings after parsing
args = None


def resource_path(*path_parts):
    """
    Return absolute path to a resource.
    - In PyInstaller: path is inside _MEIPASS/shareranger/...
    - In source mode: path is inside project_root/rules/...
    """
    if hasattr(sys, "_MEIPASS"):
        # When bundled, files are under _MEIPASS/shareranger/
        base_path = os.path.join(sys._MEIPASS, "shareranger")
    else:
        # From source: base path is project root
        current_file_dir = os.path.abspath(
            os.path.dirname(__file__)
        )  # shareranger/config/
        base_path = os.path.abspath(os.path.join(current_file_dir, ".."))
    return os.path.join(base_path, *path_parts)


# Defaults
# Default rule directory for loading rules
DEFAULT_RULE_DIR = resource_path("rules", "rulesets")
# Max file content bytes to read for content matching
DEFAULT_MAX_BYTES = 1024 * 1024 * 5  # 5MB
# Default snippet length for match previews
DEFAULT_SNIPPET_LENGTH = 50
# Default number of threads for host scans
DEFAULT_HOST_THREADS = 100
# Default number of threads for share scans
DEFAULT_SHARE_THREADS = 30
# Default number of threads for file processing
DEFAULT_FILE_THREADS = 50
# Timestamped filenames for result files
_timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
DEFAULT_RESULT_FILE = f"{_timestamp}-ShareRanger.ndjson"
DEFAULT_RESULT_FILE_CSV = f"{_timestamp}-ShareRanger.csv"


def parse_args():
    global args
    parser = argparse.ArgumentParser(
        description="ShareRanger – Fileshare Discovery & Analysis Tool"
    )
    parser.add_argument(
        "--dc",
        type=str,
        help="Domain Controller for LDAP searches (e.g. dc.example.com)",
    )
    parser.add_argument(
        "--base-dn", type=str, help="Base DN for LDAP searches (e.g. DC=example,DC=com)"
    )
    parser.add_argument(
        "--local-path",
        type=str,
        help="Path to a local folder to scan (skips SMB host and DFS discovery)",
    )
    parser.add_argument(
        "--rule-dir",
        type=str,
        default=DEFAULT_RULE_DIR,
        help="Base directory to recursively load all rule YAMLs",
    )
    parser.add_argument(
        "--max-bytes",
        type=int,
        default=DEFAULT_MAX_BYTES,
        help="Max file content bytes to read for content matching (Default: 5MB)",
    )
    parser.add_argument(
        "--snippet-length",
        type=int,
        default=DEFAULT_SNIPPET_LENGTH,
        help="Number of characters to include as match preview (Default: 50)",
    )
    parser.add_argument(
        "--host-threads",
        type=int,
        default=DEFAULT_HOST_THREADS,
        help="Max concurrent threads for host scans (Default: 100)",
    )
    parser.add_argument(
        "--share-threads",
        type=int,
        default=DEFAULT_SHARE_THREADS,
        help="Max concurrent threads for scanning shares (Default: 30)",
    )
    parser.add_argument(
        "--file-threads",
        type=int,
        default=DEFAULT_FILE_THREADS,
        help="Max concurrent threads for processing files (Default: 50)",
    )
    parser.add_argument(
        "--export-csv", action="store_true", help="Write findings to CSV file"
    )
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Enable verbose (debug) output",
    )

    args = parser.parse_args()
