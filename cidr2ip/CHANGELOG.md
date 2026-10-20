# 3.3.3

- Built with Go 1.27.3, which includes the fix for CVE-2026-111943 (`crypto/x509`).
- Version now matches `ipinfo`; all CLIs are released together.
- Dropped the 32-bit Windows ARM (`windows_arm`) build, which Go no longer supports.
- Errors now go to stderr instead of stdout.
- Windows release zip now contains `cidr2ip.exe` instead of the versioned filename.
- Install scripts (`deb.sh`, `windows.ps1`) now fail loudly when installation does not succeed.

# 1.0.0

- `cidr2ip` converts CIDRs to individual IPs within those CIDRs.
- The command exists as a standalone binary in addition to as a subcommand on the
  main `ipinfo` command.
- `cidr2ip` accepts inputs the same way our other commands do, via stdin, args,
  files, etc.
