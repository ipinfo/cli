# 3.3.3

- Built with Go 1.27.3, which includes the fix for CVE-2026-111943 (`crypto/x509`).
- Version now matches `ipinfo`; all CLIs are released together.
- Dropped the 32-bit Windows ARM (`windows_arm`) build, which Go no longer supports.
- Errors now go to stderr instead of stdout.
- Windows release zip now contains `randip.exe` instead of the versioned filename.
- Install scripts (`deb.sh`, `windows.ps1`) now fail loudly when installation does not succeed.

# 1.1.0

- `randip now supports` allowing generation of only unique IPs via the `--unique`
  or `-u` flags.
- Additionally, a bug was fixed where the ending IP of a range of size greater
  than 1 was never selected.
