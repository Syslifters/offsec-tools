# -*- mode: python ; coding: utf-8 -*-
# PyRDP MITM: twisted and other deps are loaded dynamically; collect_all is required.

from PyInstaller.utils.hooks import collect_all, collect_submodules

block_cipher = None

datas = [('pyrdp/mitm/mitm.default.ini', 'pyrdp/mitm')]
binaries = []
hiddenimports = [
    'twisted.internet.asyncioreactor',
    'twisted.internet.reactor',
    'twisted.plugins',
    'asyncio',
]

for pkg in (
    'twisted',
    'pyrdp',
    'cryptography',
    'scapy',
    'OpenSSL',
    'zope.interface',
    'incremental',
    'constantly',
    'automat',
    'hyperlink',
    'attrs',
    'service_identity',
    'pyasn1',
    'rsa',
    'pycryptodome',
    'Crypto',
    'numpy',
    'progressbar',
    'appdirs',
    'namesgenerator',
    'pytz',
):
    try:
        tmp = collect_all(pkg)
        datas += tmp[0]
        binaries += tmp[1]
        hiddenimports += tmp[2]
    except Exception:
        pass

hiddenimports += collect_submodules('pyrdp')

a = Analysis(
    ['pyrdp/bin/mitm.py'],
    pathex=[],
    binaries=binaries,
    datas=datas,
    hiddenimports=hiddenimports,
    hookspath=[],
    hooksconfig={},
    runtime_hooks=[],
    excludes=[],
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
    name='pyrdp-mitm',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    upx_exclude=[],
    runtime_tmpdir=None,
    console=True,
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
)
