# shareranger/utils/authz.py

from __future__ import annotations
import ctypes
import ctypes.wintypes as w
import win32security

# ––––– authz.dll prototypes ––––––––––––––––––––––––––––––––––––––––––––––
_authz = ctypes.WinDLL("authz", use_last_error=True)
AUTHZ_NO_AUDIT = 0x00000001
SE_GROUP_ENABLED = 0x00000004


class SID_AND_ATTRIBUTES(ctypes.Structure):
    _fields_ = [
        ("Sid", ctypes.c_void_p),
        ("Attributes", w.DWORD),
    ]


class LUID(ctypes.Structure):
    _fields_ = [
        ("LowPart", w.DWORD),
        ("HighPart", w.LONG),
    ]


PLUID = ctypes.POINTER(LUID)

# Resource‐Manager initialisieren / freigeben
_AuthzInitializeRM = _authz.AuthzInitializeResourceManager
_AuthzInitializeRM.restype = w.BOOL
_AuthzInitializeRM.argtypes = [
    w.DWORD,
    w.LPVOID,
    w.LPVOID,
    w.LPVOID,
    w.LPCWSTR,
    ctypes.POINTER(ctypes.c_void_p),
]
_AuthzFreeRM = _authz.AuthzFreeResourceManager
_AuthzFreeRM.argtypes = [ctypes.c_void_p]

# Context initialisieren / Gruppen hinzufügen / freigeben
_AuthzInitCtxFromSid = _authz.AuthzInitializeContextFromSid
_AuthzInitCtxFromSid.restype = w.BOOL
_AuthzInitCtxFromSid.argtypes = [
    w.DWORD,
    ctypes.c_void_p,
    ctypes.c_void_p,
    ctypes.c_void_p,
    PLUID,
    ctypes.c_void_p,
    ctypes.POINTER(ctypes.c_void_p),
]
_AuthzAddSids = _authz.AuthzAddSidsToContext
_AuthzAddSids.argtypes = [
    ctypes.c_void_p,
    ctypes.POINTER(SID_AND_ATTRIBUTES),
    w.DWORD,
    ctypes.POINTER(ctypes.c_void_p),
]
_AuthzAddSids.restype = w.BOOL

_AuthzFreeCtx = _authz.AuthzFreeContext
_AuthzFreeCtx.argtypes = [ctypes.c_void_p]


# Access‐Check-Strukturen
class AUTHZ_ACCESS_REQUEST(ctypes.Structure):
    _fields_ = [
        ("DesiredAccess", w.DWORD),
        ("PrincipalSelfSid", ctypes.c_void_p),
        ("ObjectTypeList", ctypes.c_void_p),
        ("ObjectTypeListLength", w.DWORD),
        ("OptionalArguments", ctypes.c_void_p),
    ]


class AUTHZ_ACCESS_REPLY(ctypes.Structure):
    _fields_ = [
        ("ResultListLength", w.DWORD),
        ("GrantedAccessMask", ctypes.POINTER(w.DWORD)),
        ("SaclEvaluationResults", ctypes.POINTER(w.DWORD)),
        ("Error", ctypes.POINTER(w.DWORD)),
    ]


_AuthzAccessCheck = _authz.AuthzAccessCheck
_AuthzAccessCheck.restype = w.BOOL
_AuthzAccessCheck.argtypes = [
    w.DWORD,
    ctypes.c_void_p,
    ctypes.POINTER(AUTHZ_ACCESS_REQUEST),
    ctypes.c_void_p,
    ctypes.c_void_p,
    ctypes.c_void_p,
    w.DWORD,
    ctypes.POINTER(AUTHZ_ACCESS_REPLY),
    ctypes.c_void_p,
]


# ––––– Hilfsfunktionen –––––––––––––––––––––––––––––––––––––––––––––––––––––
def _get_sd_ptr(path: str) -> tuple[ctypes.c_void_p, ctypes.Array]:
    """
    Hole die SECURITY_DESCRIPTOR eines Pfads und pinne die Bytes in einem
    ctypes-Buffer. Rückgabe: (Pointer auf SD, Buffer-Objekt).
    """
    sd = win32security.GetFileSecurity(
        path,
        win32security.OWNER_SECURITY_INFORMATION
        | win32security.DACL_SECURITY_INFORMATION,
    )
    buf = ctypes.create_string_buffer(bytes(sd))
    return ctypes.cast(buf, ctypes.c_void_p), buf


def _obj_to_ptr(sid_obj) -> tuple[ctypes.c_void_p, ctypes.Array]:
    """
    Konvertiere einen PySID (z.B. von ConvertStringSidToSid) in einen
    gepinnten ctypes-Buffer plus c_void_p darauf.
    """
    raw = bytes(sid_obj)
    buf = ctypes.create_string_buffer(raw)
    return ctypes.cast(buf, ctypes.c_void_p), buf


# ––––– Klasse für Authz-Kontext ––––––––––––––––––––––––––––––––––––––––––––
class AuthzContext:
    """
    AuthzContext kapselt Resource Manager (rm) und Client Context (ctx).
    Die Puffer, die SID-Bytes und SID_AND_ATTRIBUTES-Arrays enthalten,
    werden in _keep gehalten, damit sie nicht vorzeitig GC’ed werden.

    Beispiel:
        with AuthzContext(user_sid, extra_sids) as authz:
            grant = authz.check_access(sd_ptr, MAXIMUM_ALLOWED)
    """

    def __init__(self, user_sid: str, extra_sids: set[str] | None = None):
        self._keep: list[ctypes.Array] = []
        self._rm = ctypes.c_void_p()
        self._ctx = ctypes.c_void_p()

        # 1) Resource Manager initialisieren
        if not _AuthzInitializeRM(
            AUTHZ_NO_AUDIT, None, None, None, "ShareRanger", ctypes.byref(self._rm)
        ):
            raise ctypes.WinError(ctypes.get_last_error())

        # 2) Basis-Kontext vom Benutzer-SID
        user_sid_obj = win32security.ConvertStringSidToSid(user_sid)
        ptr, buf = _obj_to_ptr(user_sid_obj)
        self._keep.append(buf)

        if not _AuthzInitCtxFromSid(
            0, ptr, self._rm, None, None, None, ctypes.byref(self._ctx)
        ):
            _AuthzFreeRM(self._rm)
            raise ctypes.WinError(ctypes.get_last_error())

        # 3) Extra-SIDs hinzufügen (falls vorhanden)
        if extra_sids:
            sid_attrs: list[SID_AND_ATTRIBUTES] = []
            for s in extra_sids:
                sid_obj = win32security.ConvertStringSidToSid(s)
                sid_ptr, sid_buf = _obj_to_ptr(sid_obj)
                self._keep.append(sid_buf)
                sid_attrs.append(SID_AND_ATTRIBUTES(sid_ptr, SE_GROUP_ENABLED))

            arr = (SID_AND_ATTRIBUTES * len(sid_attrs))(*sid_attrs)
            self._keep.append(arr)  # pinne das Array

            new_ctx = ctypes.c_void_p()
            if not _AuthzAddSids(self._ctx, arr, len(arr), ctypes.byref(new_ctx)):
                _AuthzFreeCtx(self._ctx)
                _AuthzFreeRM(self._rm)
                raise ctypes.WinError(ctypes.get_last_error())

            _AuthzFreeCtx(self._ctx)
            self._ctx = new_ctx

    def check_access(self, path: str, desired_mask: int) -> int:
        """
        Führt AuthzAccessCheck durch. Gibt den granted-Mask zurück.
        """
        sd_ptr, buf = _get_sd_ptr(path)

        req = AUTHZ_ACCESS_REQUEST()
        req.DesiredAccess = desired_mask

        granted = w.DWORD()
        sacl = w.DWORD()
        error = w.DWORD()
        reply = AUTHZ_ACCESS_REPLY(
            1,
            ctypes.pointer(granted),
            ctypes.pointer(sacl),
            ctypes.pointer(error),
        )
        if not _AuthzAccessCheck(
            0,  # Flags
            self._ctx,
            ctypes.byref(req),
            None,  # AuditEvent
            sd_ptr,
            None,
            0,  # keine zusätzlichen SDs
            ctypes.byref(reply),
            None,  # Ergebnis-Handle
        ):
            raise ctypes.WinError(ctypes.get_last_error())

        return granted.value

    def free(self):
        """Räumt Context und Resource Manager auf."""
        if getattr(self, "_ctx", None):
            _AuthzFreeCtx(self._ctx)
            self._ctx = None
        if getattr(self, "_rm", None):
            _AuthzFreeRM(self._rm)
            self._rm = None
        self._keep.clear()

    def __enter__(self) -> AuthzContext:
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.free()

    def __del__(self):
        self.free()
