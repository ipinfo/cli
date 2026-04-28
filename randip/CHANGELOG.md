# 3.3.2

- Errors now go to stderr instead of stdout.
- Windows release zip now contains `randip.exe` instead of the versioned filename.
- Install scripts (`deb.sh`, `windows.ps1`) now fail loudly when installation does not succeed.

# 1.1.0

- `randip now supports` allowing generation of only unique IPs via the `--unique`
  or `-u` flags.
- Additionally, a bug was fixed where the ending IP of a range of size greater
  than 1 was never selected.
