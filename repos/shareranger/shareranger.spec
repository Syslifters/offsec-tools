# -*- mode: python ; coding: utf-8 -*-

import sys
from PyInstaller.utils.hooks import collect_submodules, collect_dynamic_libs  # +++

# Bundle all yaml rules
datas = [
    ('shareranger/rules/rulesets/*.yaml', 'shareranger/rules/rulesets')
]

# Hidden imports
hiddenimports = [
    'httpx','httpcore','idna','anyio','trio','trio.socket','sniffio',
    'ldap3','pyasn1','pyasn1.compat.octets','collections.abc','UserDict',
    'dns.asyncquery','dns._asyncio_backend','dns._trio_backend',
    'dns.quic','dns.quic._trio','dns.quic._common','dns.quic._asyncio','dns.quic._sync',
    'aioquic','aioquic.quic','aioquic.h3',
    'win32com','win32com.client','win32com.gen_py',
    'bcrypt','future','backports.ssl_match_hostname',
    'winkerberos'  # +++ ensure SSPI backend is present
] + collect_submodules("win32com")

# Collect native binaries for winkerberos
binaries = []
excludes = []

if sys.platform == "win32":
    binaries += collect_dynamic_libs('winkerberos')  # +++
    # Avoid MIT-GSSAPI on Windows so it won't look for KfW
    excludes += ['gssapi', 'gssapi.raw']             # +++

block_cipher = None

a = Analysis(
    ['build_shareranger.py'],
    pathex=[],
    binaries=binaries,       # +++
    datas=datas,
    hiddenimports=hiddenimports,
    hookspath=[],
    hooksconfig={},
    runtime_hooks=[],
    excludes=excludes,       # +++
    win_no_prefer_redirects=False,
    win_private_assemblies=False,
    cipher=block_cipher,
    noarchive=False,
)

pyz = PYZ(a.pure, a.zipped_data, cipher=block_cipher)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.zipfiles,
    a.datas,
    [],
    name='shareranger',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    upx_exclude=[],
    runtime_tmpdir=None,
    console=True,
)
