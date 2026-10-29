# 3.3.4

- Fixed the PowerShell install script (`windows.ps1`) failing after extracting the binary.

# 3.3.3

- Built with Go 1.27.3, which includes the fix for CVE-2026-111943 (`crypto/x509`).
- Version now matches `ipinfo`; all CLIs are released together.
- Dropped the 32-bit Windows ARM (`windows_arm`) build, which Go no longer supports.
- Errors now go to stderr instead of stdout.
- Fixed `--include-cidrs` / `--include-ranges` dropping every CIDR and range match when combined with `--exclude-reserved` (`-x`).
- Windows release zip now contains `grepip.exe` instead of the versioned filename.
- Install scripts (`deb.sh`, `windows.ps1`) now fail loudly when installation does not succeed.
- Added new flags: `--include-cidrs`, `--include-ranges`, `--cidrs-only`, `--ranges-only`

# 1.2.2

- Some special bogon edge cases are now properly ignored by grepip when requested.
