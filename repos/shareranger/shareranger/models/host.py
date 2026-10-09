from dataclasses import dataclass, field, asdict
from typing import List
from shareranger.models.share import Share


@dataclass
class Host:
    name: str
    online: bool = False
    shares: List[Share] = field(default_factory=list)

    def add_share(self, share: Share) -> None:
        share.host_name = self.name
        self.shares.append(share)

    def to_dict(self):
        return {
            "name": self.name,
            "online": self.online,
            "shares": [s.to_dict() for s in self.shares],
        }

    @staticmethod
    def from_dict(data):
        shares = [Share(**s) for s in data.get("shares", [])]
        return Host(name=data["name"], online=data.get("online", False), shares=shares)
