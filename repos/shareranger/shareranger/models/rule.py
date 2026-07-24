from dataclasses import dataclass, field
from typing import List, Optional, Dict, Any
from enum import Enum
from shareranger.logging.logger import log_warning


class RuleType(str, Enum):
    FILE = "File"
    DIRECTORY = "Directory"
    SHARE = "Share"


class RuleAction(str, Enum):
    KEEP = "Keep"
    DISCARD = "Discard"
    MATCH = "Match"


class RuleLocation(str, Enum):
    FILE_PATH = "FilePath"
    FILE_NAME = "FileName"
    FILE_EXTENSION = "FileExtension"
    SHARE_NAME = "ShareName"
    FILE_CONTENT = "FileContent"


class RuleSeverity(str, Enum):
    CRITICAL = "Critical"
    HIGH = "High"
    MEDIUM = "Medium"
    LOW = "Low"
    INFO = "Info"


@dataclass
class Rule:
    id: str
    type: RuleType
    description: str
    action: RuleAction
    location: RuleLocation
    severity: RuleSeverity
    patterns: List[str]
    ntfs: Optional[Dict[str, Any]] = None
    relay: Optional[List[str]] = None
    children: List["Rule"] = field(default_factory=list)

    def to_dict(self) -> Dict[str, Any]:
        return {
            "id": self.id,
            "type": self.type.value,
            "description": self.description,
            "action": self.action.value,
            "location": self.location.value,
            "severity": self.severity.value,
            "patterns": self.patterns,
            "ntfs": self.ntfs,
            "relay": self.relay,
        }

    @staticmethod
    def from_dict(data: Dict[str, Any], source: str) -> Optional["Rule"]:
        id = data.get("id")
        if not id:
            log_warning(f"Skipping rule with missing 'id' in {source}")
            return None

        # Case-insensitive Enum lookup
        def ci_enum(enum_cls, value):
            if not isinstance(value, str):
                raise ValueError(f"{value} is not a string")
            for member in enum_cls:
                if member.value.lower() == value.lower():
                    return member
            raise ValueError(f"Invalid value '{value}' for enum {enum_cls.__name__}")

        try:
            type_enum = ci_enum(RuleType, data["type"])
            action_enum = ci_enum(RuleAction, data["action"])
            location_enum = ci_enum(RuleLocation, data["location"])
            severity_enum = ci_enum(RuleSeverity, data["severity"])
        except Exception as e:
            log_warning(
                f"Skipping rule '{id}' due to invalid enum field in {source}: {e}"
            )
            return None

        if not isinstance(data.get("description"), str):
            log_warning(
                f"Skipping rule '{id}' with invalid or missing description in {source}"
            )
            return None

        patterns = data.get("patterns")
        if not isinstance(patterns, list) or not patterns:
            log_warning(
                f"Skipping rule '{id}' with invalid or empty patterns in {source}"
            )
            return None

        ntfs = data.get("ntfs")
        if ntfs:
            if not any(k in ntfs for k in ("readable", "writable")):
                log_warning(
                    f"Skipping rule '{id}' in {source}: 'ntfs' must contain at least one of 'readable' or 'writable'"
                )
                return None

        relay = data.get("relay")
        if action_enum == RuleAction.MATCH:
            if not isinstance(relay, list) or not all(
                isinstance(r, str) for r in relay
            ):
                log_warning(
                    f"Skipping relay rule '{id}' with invalid relay targets in {source}"
                )
                return None

        return Rule(
            id=id,
            type=type_enum,
            description=data.get("description"),
            action=action_enum,
            location=location_enum,
            severity=severity_enum,
            patterns=patterns,
            ntfs=ntfs,
            relay=relay,
        )
