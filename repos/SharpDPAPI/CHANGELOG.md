# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [1.12.0] - 2024-12-14

### Added
    * Support for local/domain credkey usage for masterkey decryption
    * Slack support to statekeys/cookies commands (@Lee Christensen)
    * Ability  to specify local state file (@Lee Christensen)
    * RPC (MS-BKUP) masterkey retrieval (@Kiblyn11)
    * User-context unprotect option for certificates (@ptr0x1)
    * Teams statekey support (@fsacer)
    * Ability to dump hashes in jtr/hashcat format (@rxwx)
    * HMAC validation to 3DES SHA1 (@rxwx)
    * SCCM object data parsing (@guervild)
    * SCCM to machine triage (@subat0mik)

### Fixed
    * User local admin SID enumeration (@clod81)
    * FIPS compliant API usage (@Calvin Hedle)
    * Updated for latest editthiscookie format (@djhohnstein)
    * String types + code cleanup (@Lee Christensen)
    * SCCM regex fix (@guervild)
    * Bug in user SID extraction (@rxwx)
    * nameof to true name for BCrypt* defs (@djhohnstein)

### Changed
    * NTLM hash now manually specified with `/ntlm` instead of inferred from `/password`
    * Backupkey now not line wrapped by default


## [1.11.3] - 2022-06-06

### Fixed
* **SharpDPAPI** project
    * `masterkeys` command now accepts a `/password:X` argument with `/target:X`


## [1.11.2] - 2022-01-12

### Fixed
* **SharpChrome** project
    * Chrome cookie file location updated for newer versions


## [1.11.1] - 2021-03-05

### Fixed
* **SharpChrome** project
    * Chrome statekey usage bug when triaging multiple users

### Changed
* **SharpChrome** project
    * Default cookie/logins/statekeys triage behavior is now to triage the current user if elevated, unless pvk/password/masterkeys specified
    * `/target:X` can now be a `C:\Users\USER\` folder for any specified triage
    * Added Brave statekey triage to `statekeys`
    * Cleaned up Chromium triage code
* Removed out of date SharpDPAPI.cna aggressor script


## [1.11.0] - 2021-03-01

### Added
* **SharpDPAPI** project
    * Added `keepass` command - ProtectedUserKey.bin decryption
    * Added `/entropy` flag to `blob` command

### Fixed
* **SharpDPAPI** project
    * Decrypted null bytes in certificate description fields messing up output


## [1.10.0] - 2021-02-25

### Added
* **SharpDPAPI** project
    * CNG private key decryption support \m/
    * Additional CAPI/CNG cert search locations
    * `/nowrap` flag to the `backupkey` command

### Fixed
* **SharpDPAPI** project
    * Bug where some extracted key components ending in 00 caused error cases

### Changed
* **SharpDPAPI** project
    * Only decrypted private keys with certs present displayed by default, the `/showall` flag for `certificates` will display all decrypted results
    * Combined `machinecerts` into `certificates /machine`
    * Corrected method for SHA1 MS hash computation fuckery with entropy (thanks @gentilkiwi)
    * Re-added certificate triage to `triage` and `machinetriage`


## [1.9.2] - 2021-01-04

### Added
* **SharpDPAPI** project
    * /target option for machinecertificates
    * more certificate information on extraction (including Enhanced Key Usages)

### Fixed
* **SharpDPAPI** project
    * User certificate extraction corrected
    * Few formatting issues


## [1.9.1] - 2020-11-05

### Added
* **SharpDPAPI** project
    * Ability to triage masterkey targets (or folder of targets) manually

### Added
* **SharpChrome** project
    * Added Chromium-based brave support
    * Added `/quiet` flag for csv output

### Fixed
* **SharpChrome** project
    * Filtering fixes for cookies


## [1.8.0] - 2020-07-13

### Added
* **SharpDPAPI** project
    * Landed @leechristensen's `search` command to search for DPAPI blobs

### Removed
* **SharpDPAPI** project
    * Removed machine/user certificate triage from the `triage` and `machinetriage` commands

### Changed
* Code cleanup and refactoring


## [1.7.0] - 2020-05-06

### Added
* **SharpDPAPI** project
    * Landed @leftp's `certificates` and `machinecerts` commands
    * Added `certificates` and `machinecerts` entries to the README.md
    * Added certificate triage to the `triage` and `machinetriage` commands
    * Using /password:X now causes the DPAPI masterkey cache to be output
*  **SharpChrome** project:
    * Using /password:X now causes the DPAPI masterkey cache to be output


## [1.6.1] - 2020-03-29

### Changed
* **SharpChrome** project
    * '/password:X' integration
* **SharpDPAPI** project
    * Combined TriageUserMasterKeysWithPass into TriageUserMasterKeys
    * '/password:X' now properly works in SharpDPAPI while elevated, as well as remotely


## [1.6.0] - 2020-03-27

### Added
* **SharpChrome** project
    * Integrated new Chrome (v80+) AES statekey decryption from @djhohnstein's SharpChrome project.
* **SharpDPAPI** project
    * landed @lefterispan's PR that incorporates plaintext password masterkey decryption.
    * expanded the PR to allow /password specification for all SharpDPAPI functions


## [1.5.1] - 2019-12-18

### Added
* **SharpChrome** project
    * **/setneverexpire** flag for **/format:json** output for **cookies** that sets the expiration date to now + 100 years

### Fixed
* **SharpChrome** project
    * Cookie datetime value parsing to prevent error conditions on invalid input.

### Changed
* **SharpChrome** project
    * Some file path output.


## [1.5.0] - 2019-07-25

### Added
* **ps** command to decrypt exported PSCredential xmls (thanks for the idea @gentilkiwi ;)
* **blob** section for the README

### Changed
* **blob** command outputs hex if the data doesn't appear to be text


## [1.4.0] - 2019-05-22

### Added
* **SharpChrome** project
    * Separate project that implements a SQLite parsing database for Chrome triage. Uses shared files with SharpDPAPI. Adapted from the SharpWeb/SharpChrome project by @djhohnstein.
    * **logins** function
        * Finds/decrypts Chrome 'Login Data' files. See README.md for complete syntax/flags.
    * **cookies** function
        * Finds/decrypts Chrome 'Cookies' files. See README.md for complete syntax/flags.
* Added **/mkfile:FILE** argument to credentials/vaults/rdg/triage commands, takes a SharpDPAPI or Mimikatz formatted file of {GUID}:SHA1 masterkey mappings (for offline triage)

### Changed
* Cleaned up and simplified the credentials/vaults/rdg/triage command functions in SharpDPAPI
* Cleaned up and reorganized SharpDPAPI's default help menu output


## [1.3.1] - 2019-05-09

### Changed
* When using /server:X, .RDG files parsed from RDCMan.settings files are translated to \\\\UNC paths for parsing

### Fixed
* **triage** command when used against a remote /server:X now works properly


## [1.3.0] - 2019-05-09

### Added
* **rdg** action
    * Find RDCMan.settings and linked .RDG files, or take a given /target .RDG/RDCMan.settings file/folder, and decrypt passwords given a /pvk, GUID key lookup table, or CryptUnprotectData (with /unprotect).
* **blob** action
    * Describe a supplied DPAPI binary blob, optionally decryption the blob with masterkey GUID lookups or a PVK masterkey decryption

### Changed
* Added IsTextUnicode() for vault/credential/blob decryption display, showing hex if unicode is detected
* Added /target:C:\FOLDER\ option for the **masterkeys** function, for offline masterkey decryption
* Updated README


## [1.2.0] - 2019-03-24 (Troopers edition ;)

### Added
* **masterkeys/vaults/creds/triage** actions
    * Remote server support for user vault/credential triage with /server:X
* **machinemasterkeys** perform master key triage for the local machine
    * implicitly elevates to SYSTEM to extract the machine's local DPAPI key
    * uses this key to triage all machine Credential files
* **machinecredentials** perform Credential file triage for the local machine
    * implicitly elevates to SYSTEM via the **machinemasterkeys** approach
    * uses the extracted masterkeys to decrypt any Credential files
* **machinevaults** perform vault triage for the local machine
    * implicitly elevates to SYSTEM via the **machinemasterkeys** approach
    * uses the extracted masterkeys to decrypt any machine Vaults
* **machinetriage** performs all machine triage actions (currently vault and credential)
    * implicitly elevates to SYSTEM via the **machinemasterkeys** approach

### Changed
* Expanded Vault credential format to handle vault credential clear attributes
* Expanded machine vault/credential search locations
* Broke out commands/files into the same general structure as Rubeus


## [1.1.1] - 2019-03-15

### Added
* **SharpDPAPI.cna** Cobalt Strike aggressor script to automate the usage of SharpDPAPI (from @leechristensen)

### Changed
* Wrapped main in try/catch

### Fixed
* Fixed Policy.vpol parsing to handle the "KSSM" (?) format. Thank you @gentilkiwi :)


## [1.1.0] - 2019-03-14

### Added
* **masterkeys** action
    * decrypts currently reachable master keys (current users or all if elevated) and attempts to decrypt them using a passed {GUI}:SHA1 masterkey lookup table, or a /pvk base64 blob representation of the domain DPAPI backup key
* **credentials** action
    * decrypts currently reachable Credential files (current users or all if elevated) and attempts to decrypt them using a passed {GUI}:SHA1 masterkey lookup table, or a /pvk base64 blob representation of the domain DPAPI backup key
* **vaults** action
    * decrypts currently reachable Vault files (current users or all if elevated) and attempts to decrypt them using a passed {GUI}:SHA1 masterkey lookup table, or a /pvk base64 blob representation of the domain DPAPI backup key
* **triage** action
    * performs all triage actions (currently vault and credential)
* CHANGELOG

### Changed
* modified the argument formats for the **backupkey** command
* retructured files so code isn't in a single file
* revamped README


## [1.0.0] - 2018-08-22

* Initial release
