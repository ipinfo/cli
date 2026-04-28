# 3.3.2

- Errors now go to stderr instead of stdout.
- Fixed `--include-cidrs` / `--include-ranges` dropping every CIDR and range match when combined with `--exclude-reserved` (`-x`).
- Windows release zip now contains `grepip.exe` instead of the versioned filename.
- Install scripts (`deb.sh`, `windows.ps1`) now fail loudly when installation does not succeed.

# 1.2.3

- Added new flags: `--include-cidrs`, `--include-ranges`, `--cidrs-only`, `--ranges-only`

# 1.2.2

- Some special bogon edge cases are now properly ignored by grepip when requested.
