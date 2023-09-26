## cidr2range-1.2.0
* When `cidr2range` accepts a file, it now also looks to see if there is a
header in CSV form, and if so changes the first column to range, just as it
changes the first non-header columns from CIDRs to IP ranges.
