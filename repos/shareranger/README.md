
# 🔍 ShareRanger

**ShareRanger** is a powerful SMB share scanner for exploring Windows network shares in enterprise environments. It discovers hosts via Active Directory and applies customizable rules to find accessible or sensitive files, directories, and shares—ideal for both offensive and defensive use cases.

## 🚀 Features

- Active Directory integration to discover SMB hosts and DFS shares
- Rule-based engine for shares, directories, and files
- YAML-based rule definitions with relay/nesting support
- File content inspection with pattern matching 
- Multi-threaded scanning (host, share, and file level)
- Timestamped JSON and CSV output

## 🧪 Usage

### Basic Scan (Discovery + Rule Matching)

Running ShareRanger with no arguments will:

- **Phase 1 – Discovery**: 
  - Finds DFS namespaces and resolves physical DFS targets
  - Queries AD for computer objects and enumerates their SMB shares
- **Phase 2 – Filesystem Walk**:
  - Walks discovered shares and applies rule-matching to:
    - Directories (structure-based rules)
    - Files (name, extension, or content-based rules)
    - File-level rule matching is parallelized via threads for performance
- **Export Scan Results**:
  - Export results to timestamped files

```bash
python -m shareranger.cli
```

---

### Export findings to CSV

Add `--export-csv` to generate a CSV with the findings alongside the JSON scan result output:

```bash
python -m shareranger.cli --export-csv
```

---

### Custom Rule Sets

Override the default rule directory:

```bash
python -m shareranger.cli --rule-dir custom_rules/
```

---

### Increase Performance with Threads

Use threading to improve performance on large networks:

```bash
--host-threads <N>   # Concurrent host scans (default: 100)
--share-threads <N>  # Concurrent share walks (default: 30)
--file-threads <N>   # Concurrent file rule checks per share (default: 50)
```

---

### Optional Manual AD Configuration

Specify domain controller and computer search base manually, instead of automatic detection:

```bash
--dc DC1.company.local --base-dn "DC=company,DC=local"
```

## 🧱 Project Structure

```
shareranger/
├── cli.py                   # Main entry point
├── discovery/               # AD & SMB scanning (AD, DFS, shares)
│   ├── ad.py
│   └── shares.py
├── models/                  # Data models: Host, Share, Rule, Match
├── rules/                   # Rule loading and pattern matching
│   ├── loader.py
│   └── matcher.py
│   └── rulesets/            # YAML rule definitions
├── utils/                   # Logger, scan result exporter
├── walker/                  # Filesystem walking logic
│   └── fswalker.py          # Multi-threaded directory & file scanning
└── config/                  # Global configuration options
    └── config.py
```

## ✅ Requirements

- Python 3.9+
- Windows OS (for Win32 + AD support)

## 🛠 Installation & Executable Build

### Running from Source

If you're using ShareRanger via source:

1. Clone the repo
2. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```
3. Run it:
   ```bash
   python -m shareranger.cli
   ```

### Building the Standalone Executable (Windows Only)

You can generate a standalone `.exe` using [PyInstaller](https://pyinstaller.org/). This bundles the entire tool (including rules and dependencies) into a single executable that works without Python installed.

#### Build Instructions

1. Install PyInstaller:
   ```bash
   pip install pyinstaller
   ```

2. Build the executable:
   ```bash
   pyinstaller shareranger.spec
   ```

3. The output executable will be in the `dist/` directory:
   ```
   dist/
   └── shareranger.exe
   ```

4. You can now run `shareranger.exe` on any compatible Windows machine — no Python required.