# 3.3.2

- Errors now go to stderr instead of stdout.
- Windows release zip now contains `prips.exe` instead of the versioned filename.
- Install scripts (`deb.sh`, `windows.ps1`) now fail loudly when installation does not succeed.

# 1.0.0

- `prips` is a command for printing all the IPs that exist within a CIDR and/or
  IP range; it is effectively a combination of `range2ip` and `cidr2ip` in one.
- The command exists as a standalone binary in addition to as a subcommand on the
  main `ipinfo` command.
- `prips` accepts inputs the same way our other commands do, via stdin, args,
  files, etc.
