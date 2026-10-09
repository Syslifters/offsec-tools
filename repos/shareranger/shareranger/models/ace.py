from dataclasses import dataclass, field
from functools import cached_property
from typing import List, Dict, Any
import ntsecuritycon
import win32security, ntsecuritycon

from shareranger.authz.authz import AuthzContext
from shareranger.logging.logger import log_debug


@dataclass
class ACE:
    path: str
    sid: str
    allow_mask: int = 0
    deny_mask: int = 0

    MAXIMUM_ALLOWED_MASK = 0x02000000

    READ_FLAGS = (
        ntsecuritycon.FILE_LIST_DIRECTORY
        | ntsecuritycon.FILE_READ_DATA
        | ntsecuritycon.FILE_READ_ATTRIBUTES
        | ntsecuritycon.FILE_READ_EA
        | ntsecuritycon.READ_CONTROL
        | ntsecuritycon.WRITE_DAC  # for adding read permissions
        | ntsecuritycon.WRITE_OWNER  # for taking ownership
    )

    WRITE_FLAGS = (
        ntsecuritycon.FILE_WRITE_DATA
        | ntsecuritycon.FILE_APPEND_DATA
        | ntsecuritycon.FILE_WRITE_ATTRIBUTES
        | ntsecuritycon.FILE_WRITE_EA
        | ntsecuritycon.FILE_DELETE_CHILD
        | ntsecuritycon.DELETE
        | ntsecuritycon.WRITE_DAC
        | ntsecuritycon.WRITE_OWNER
    )

    RISKY_RIDS = {
        "-513",  # Domain Users
        "-515",  # Domain Computers
    }

    RISKY_SIDS = {
        "S-1-1-0",  # Everyone
        "S-1-5-11",  # Authenticated Users
        "S-1-5-32-545",  # BUILTIN\Users
    }

    @staticmethod
    def decode_advanced_permissions(mask: int) -> List[str]:
        perms = []

        if mask & ntsecuritycon.FILE_EXECUTE:
            perms.append("Traverse folder / execute file")
        if (
            mask & ntsecuritycon.FILE_LIST_DIRECTORY
            or mask & ntsecuritycon.FILE_READ_DATA
        ):
            perms.append("List folder / read data")
        if mask & ntsecuritycon.FILE_READ_ATTRIBUTES:
            perms.append("Read attributes")
        if mask & ntsecuritycon.FILE_READ_EA:
            perms.append("Read extended attributes")
        if mask & ntsecuritycon.FILE_WRITE_DATA:
            perms.append("Create files / write data")
        if mask & ntsecuritycon.FILE_APPEND_DATA:
            perms.append("Create folders / append data")
        if mask & ntsecuritycon.FILE_WRITE_ATTRIBUTES:
            perms.append("Write attributes")
        if mask & ntsecuritycon.FILE_WRITE_EA:
            perms.append("Write extended attributes")
        if mask & ntsecuritycon.FILE_DELETE_CHILD:
            perms.append("Delete subfolders and files")
        if mask & ntsecuritycon.DELETE:
            perms.append("Delete")
        if mask & ntsecuritycon.READ_CONTROL:
            perms.append("Read permissions")
        if mask & ntsecuritycon.WRITE_DAC:
            perms.append("Change permissions")
        if mask & ntsecuritycon.WRITE_OWNER:
            perms.append("Take ownership")

        return perms

    @cached_property
    def permissions(self) -> List[str]:
        return ACE.decode_advanced_permissions(self.access_mask)

    @cached_property
    def security_principal(self) -> str:
        try:
            sid_obj = win32security.ConvertStringSidToSid(self.sid)
            name, domain, type = win32security.LookupAccountSid(None, sid_obj)
            return f"{domain}\\{name}" if domain else name
        except Exception as e:
            log_debug(f"Failed to resolve SID {self.sid}: {e}")
            return self.sid

    @cached_property
    def access_mask(self) -> int:
        with AuthzContext(self.sid, None) as authz:
            granted = authz.check_access(self.path, ACE.MAXIMUM_ALLOWED_MASK)

        return granted

    @cached_property
    def has_read(self) -> bool:
        return bool(self.access_mask & ACE.READ_FLAGS)

    @cached_property
    def has_write(self) -> bool:
        return bool(self.access_mask & ACE.WRITE_FLAGS)

    def to_dict(self) -> Dict[str, Any]:
        return {
            "sid": self.sid,
            "security_principal": self.security_principal,
            "access_mask": hex(self.access_mask),
            "permissions": sorted(self.permissions),
        }
