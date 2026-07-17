from dataclasses import dataclass
from typing import Optional, List, Dict, Any
from shareranger.models.ace import ACE
from shareranger.models.rule import Rule


@dataclass
class Match:
    rule: Rule
    unc_path: Optional[str]
    pattern: str
    snippet: str
    insecure_aces: Optional[List[ACE]] = None

    def to_dict(self) -> Dict[str, Any]:
        return {
            "id": self.rule.id,
            "unc_path": self.unc_path,
            "pattern": self.pattern,
            "snippet": self.snippet,
            "severity": self.rule.severity.value,
            "location": self.rule.location.value,
            "action": self.rule.action.value,
            "type": self.rule.type.value,
            "ntfs": self.rule.ntfs,
            "insecure_aces": (
                [ace.to_dict() for ace in self.insecure_aces]
                if self.insecure_aces
                else None
            ),
        }
