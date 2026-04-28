# 3.3.2

- Errors now go to stderr instead of stdout.
- Windows release zip now contains `matchip.exe` instead of the versioned filename.
- Install scripts (`deb.sh`, `windows.ps1`) now fail loudly when installation does not succeed.

# 1.0.0

- `matchip` is a command for printing the IPs and subnets that fall under the given list of subnets.
- The command exists as a standalone binary in addition to as a subcommand on the main `ipinfo` command.
- `matchip` accepts inputs the same way our other commands do, via stdin, args, files, etc.
