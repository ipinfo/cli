# 3.3.2

- Errors now go to stderr instead of stdout.
- Windows release zip now contains `range2ip.exe` instead of the versioned filename.
- Install scripts (`deb.sh`, `windows.ps1`) now fail loudly when installation does not succeed.

# 1.0.0

- `range2ip` converts IP ranges to individual IPs within those ranges.
- The command exists as a standalone binary in addition to as a subcommand on the
  main `ipinfo` command.
- `range2ip` accepts inputs the same way our other commands do, via stdin, args,
  files, etc.
