import logging
import colorama
import sys

# Enable ANSI code support in Windows terminals
sys.stdout = colorama.AnsiToWin32(sys.stdout).stream
sys.stderr = colorama.AnsiToWin32(sys.stderr).stream
colorama.init(autoreset=True)

# ANSI background color codes for severity label
COLOR_MAP = {
    "CRITICAL": "\033[41m\033[97m",  # Black background
    "HIGH": "\033[91m",  # Red background
    "MEDIUM": "\033[93m",  # Yellow background
    "LOW": "\033[92m",  # Green background
    "INFO": "\033[94m",  # Blue background
}
RESET = "\033[0m"


def init_logging(verbose: bool = False):
    """
    Initializes logging system. Should be called once at app start.
    """
    logging.basicConfig(
        level=logging.DEBUG if verbose else logging.INFO,
        format="[%(levelname)s] %(asctime)s %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )


def log_status(message):
    logging.info(f"[*] {message}")


def log_success(message):
    logging.info(f"[+] {message}")


def log_warning(message):
    logging.warning(f"[!] {message}")


def log_error(message):
    logging.error(f"[x] {message}")


def log_debug(message):
    logging.debug(f"[~] {message}")


def log_match(
    action, type, path, id, severity, location, pattern, match_snippet, verbose=False
):

    label_action = action.upper()
    label_severity = severity.upper()
    level = logging.INFO if not verbose else logging.DEBUG

    # Color only the severity label
    color = COLOR_MAP.get(label_severity, "")
    colored_severity = f"{color}[{label_severity}]{RESET}"

    logging.log(
        level,
        f"[{label_action}] {colored_severity} {type}: {path} "
        f"(rule: {id}, location: {location}, "
        f'pattern: "{pattern}", match: "{match_snippet}")',
    )
