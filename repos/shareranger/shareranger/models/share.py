from dataclasses import dataclass, field
from functools import cached_property
import os
from typing import Any, List, Dict, Optional
from shareranger.models.ace import ACE
from shareranger.models.match import Match
from shareranger.logging.logger import log_debug


@dataclass
class Share:
    _host_name: Optional[str] = field(
        default=None, init=False, repr=False, compare=False
    )
    name: str
    unc_path: str
    matches: List[Match] = field(default_factory=list)

    @property
    def host_name(self) -> Optional[str]:
        return self._host_name

    @host_name.setter
    def host_name(self, value: Optional[str]) -> None:
        self._host_name = value

    @cached_property
    def accessible(self) -> bool:
        try:
            os.listdir(self.unc_path)
            return True
        except Exception as e:
            log_debug(f"Failed to access share: {e}")

        return False

    def to_dict(self) -> Dict[str, Any]:
        return {
            "name": self.name,
            "unc_path": self.unc_path,
            "accessible": self.accessible,
            "matches": [m.to_dict() for m in self.matches],
        }
