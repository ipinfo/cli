# 3.3.2

- Errors now go to stderr instead of stdout.
- Windows release zip now contains `cidr2range.exe` instead of the versioned filename.
- Install scripts (`deb.sh`, `windows.ps1`) now fail loudly when installation does not succeed.

# 1.2.0

- When `cidr2range` accepts a file, it now also looks to see if there is a
  header in CSV form, and if so changes the first column to range, just as it
  changes the first non-header columns from CIDRs to IP ranges.
