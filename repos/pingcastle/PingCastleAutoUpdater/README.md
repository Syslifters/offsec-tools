# PingCastle Auto Updater

PingCastleAutoUpdater is an automated update tool that downloads and installs the latest versions of PingCastle from GitHub releases. It manages configuration migrations, handles file updates, and ensures your PingCastle installation stays current with minimal manual effort.

## Features

- **Automatic Version Detection**: Compares current PingCastle version with the latest GitHub release
- **Smart Configuration Handling**: Preserves existing settings when updating, with automatic XML to JSON conversion when needed
- **Initial State Migration**: Automatically migrates legacy XML configurations to JSON format before updates
- **Backup & Recovery**: Automatically backs up existing configurations before updates
- **Dry-Run Preview**: Test updates without modifying any files
- **Flexible Release Selection**: Choose stable releases, preview versions, or wait for a release stabilization period
- **Custom Update Sources**: Point to alternative update repositories if needed
- **Complete Audit Trail**: Generates detailed migration and conversion reports for all operations

## System Requirements

- **Platform**: Windows x64 only
- **.NET Runtime**: Not required - includes self-contained .NET 8.0 runtime
- **Disk Space**: Minimal (executable only, no external dependencies)

## Installation

Place `PingCastleAutoUpdater.exe` in the same directory as `pingcastle.exe`. The tool is self-contained and requires no additional software to be installed.

## Basic Usage

### Simple Update Check and Install

```bash
PingCastleAutoUpdater.exe
```

This will:
1. Check GitHub for the latest stable release
2. Compare it with your current PingCastle version
3. Download and install the update if a newer version is available
4. Handle any configuration migrations automatically

### Preview Changes with Dry-Run

```bash
PingCastleAutoUpdater.exe --dry-run
```

The `--dry-run` option is useful for:
- **Testing before applying updates**: See exactly what files would be modified
- **Verifying configuration migration**: Check if XML to JSON conversion would succeed
- **Audit purposes**: Review what changes an update would introduce
- **No risk operation**: Absolutely no files are modified in dry-run mode

In dry-run mode, you will see:
- Which files would be downloaded and extracted
- How configuration files would be merged or converted
- A detailed conversion report saved to disk for review
- Any custom settings that would be preserved

Example output:
```
[DRY-RUN] No files will be modified. Analyzing update contents...

[DRY-RUN] Would save PingCastle.exe
[DRY-RUN] Would merge .config file: PingCastle.exe.config
[DRY-RUN] Would save PingCastleCommon.dll
...
[DRY-RUN MODE - NO CHANGES WILL BE MADE]

Analyzing configuration conversion...
[OK] Analysis completed - Configuration would be converted successfully

[Report] Preview report saved: ConversionPreview_20240115_143022.txt
[Data] Sections to convert: 5
[Settings] Settings to map: 24

To proceed with conversion, run without --dry-run flag.
```

## Command-Line Switches

### `--help`
Display the help message with all available options.

```bash
PingCastleAutoUpdater.exe --help
```

### `--dry-run`
Preview all changes that would be made without modifying any files.

**Use cases:**
- Test an update before applying it to production systems
- Verify configuration migration will succeed
- Audit what a new version would change
- Run on a test machine before deploying updates

```bash
PingCastleAutoUpdater.exe --dry-run
```

**Note:** Even in dry-run mode, temporary preview reports are generated so you can review the changes in detail. These reports are automatically cleaned up after the operation completes.

### `--use-preview`
Include preview/prerelease versions when checking for updates. By default, only stable releases are considered.

```bash
PingCastleAutoUpdater.exe --use-preview
```

This is useful for:
- Testing new features before official release
- Getting early access to bug fixes
- Running on non-production systems where stability is less critical

### `--force-download`
Download and install the latest release even if your current version is already up to date. Useful for re-installing or downgrading.

```bash
PingCastleAutoUpdater.exe --force-download
```

**Warning:** This will overwrite your current installation. Consider using `--dry-run` first to see what will be replaced.

### `--wait-for-days <n>`
Only use releases that have been published for at least N days. This adds a stabilization period before updating.

```bash
PingCastleAutoUpdater.exe --wait-for-days 30
```

Use cases:
- Implement a controlled update schedule
- Allow time for community feedback on new releases
- Ensure stability in production environments
- Stagger updates across multiple installations

**Example:** With `--wait-for-days 30`, the tool will skip any release published in the last 30 days, even if it's the "latest" release on GitHub.

### `--api-url <url>`
Direct the PingCastle AutoUpdater to use PingCastle Pro or Enterprise to perform updates internally for systems that cannot connect to GitHub.

```bash
.\PingCastleAutoUpdater.exe --api-url https://pingcastle.your.server/api/release
```
**Note:** The URL must be HTTP or HTTPS. The API should return release information in the same format as GitHub's API.

### Migrating from PingCastle 3.4 to 3.5

When an existing user runs the auto-updater for the first time and has both XML and JSON configurations:

```bash
PingCastleAutoUpdater.exe
```

The tool will:
1. Detect both `PingCastle.exe.config` and `appsettings.console.json`
2. Perform initial state migration:
   - Convert XML settings to JSON format
   - Merge XML settings into the existing JSON (XML values take precedence)
   - Archive the legacy XML configuration file as `PingCastle.exe.config.bak` for potential recovery
3. Proceed with downloading and installing the latest update
4. Generate a report documenting the update's configuration changes

This ensures smooth transition for long-time users upgrading to the new JSON-based configuration system. The original XML configuration is archived and can be restored if needed.

After updates, the tool may generate several report files in the PingCastle directory:

### Migration and Conversion Reports

- `InitialStateMigration_*.txt` - Details of pre-update XML-to-JSON migration when both XML and JSON configs existed initially
  - Tracks which XML sections were converted
  - Documents which settings were merged into JSON
  - Shows which XML values took precedence over existing JSON values
  - Lists any custom XML elements that were preserved

- `ConversionReport_*.txt` - Details of post-update XML to JSON configuration conversion and merging
  - Documents any XML-to-JSON conversions during the update
  - Shows which settings were merged with update versions
  - Lists configuration changes and migrations
  - Provides audit trail for troubleshooting

## Troubleshooting

### "No current version detected"

The tool couldn't find `pingcastle.exe` in the same directory as the updater. This is normal for first-time updates. The tool will proceed with the download.

```
No current version detected - download will proceed
```

### Initial State Migration Failed

If migration of XML to JSON fails during the first run:

```
Error during initial state migration: [error details]
Restoring original configuration files...
Original JSON configuration restored.
Initial state migration error report saved...
```

**What happens:**
- Your original JSON configuration is restored from the temporary backup
- An error report is saved with details about what failed
- The update continues using your existing configuration
- You can manually review and merge settings if needed

**Next steps:**
1. Check the error report for details
2. Run with `--dry-run` to preview the migration
3. Contact support with the error report if the issue persists

### Configuration Merge Failed (Post-Update)

If configuration merging fails during the update phase:

```
Warning: Could not backup existing JSON config: [error details]
```

Your previous configuration is preserved as a `.bak` file. You may need to manually review and merge settings.

### Network Issues

If GitHub is unreachable or the API returns an error:

```
Network error: [error details]
```

Check your network connectivity. You can optionally specify a custom API URL with `--api-url` if GitHub is inaccessible from your network.

## Running as Scheduled Task

To automate updates on a Windows system, you can set up a scheduled task:

1. Open Task Scheduler
2. Create Basic Task
3. Set trigger (e.g., weekly)
4. Set action: Start `PingCastleAutoUpdater.exe` in the PingCastle directory
5. Optional: Add arguments like `--wait-for-days 7`

## Support
If you need help see [Stay Connected and Contribute](https://github.com/netwrix/pingcastle?tab=readme-ov-file#stayconnected--contribute) for how to engage in the community. If you are a paying customer you can find how to engage support from the [Features and Bugs](https://github.com/netwrix/pingcastle?tab=readme-ov-file#features--bugs) section of the repository readme.
