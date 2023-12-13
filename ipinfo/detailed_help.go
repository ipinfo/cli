package main

var DetailedHelp = `Usage: ipinfo <cmd> [<opts>] [<args>]

Commands:
  <ip>        look up details for an IP address, e.g. 8.8.8.8.
  <asn>       look up details for an ASN, e.g. AS123 or as123.
  myip        get details for your IP.
  bulk        get details for multiple IPs in bulk.
  asn         tools related to ASNs.
  summarize   get summarized data for a group of IPs.
  map         open a URL to a map showing the locations of a group of IPs.
  prips       print IP list from CIDR or range.
  grepip      grep for IPs matching criteria from any source.
  matchip     print the overlapping IPs and subnets.
  grepdomain  grep for domains matching criteria from any source.
  cidr2range  convert CIDRs to IP ranges.
  cidr2ip     convert CIDRs to individual IPs within those CIDRs.
  range2cidr  convert IP ranges to CIDRs.
  range2ip    convert IP ranges to individual IPs within those ranges.
  randip      Generates random IPs.
  splitcidr   splits a larger CIDR into smaller CIDRs.
  mmdb        read, import and export mmdb files.
  calc        evaluates a mathematical expression that may contain IP addresses.
  tool        misc. tools related to IPs, IP ranges and CIDRs.
  download    download free ipinfo database files.
  cache       manage the cache.
  config      manage the config.
  quota       print the request quota of your account.
  init        login or signup account.
  logout      delete your current API token session.
  completion  install or output shell auto-completion script.
  version     show current version.

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --nocache
      do not use the cache.
    --version, -v
      show binary release number.
    --help, -h
      show help.

  Outputs:
    --field <field>, -f <field>
      lookup only specific fields in the output.
      field names correspond to JSON keys, e.g. 'hostname' or 'company.type'.
      multiple field names must be separated by commas.
    --nocolor
      disable colored output.

  Formats:
    --pretty, -p
      output pretty format.
    --json, -j
      output JSON format.
    --csv, -c
      output CSV format.
    --yaml, -y
      output YAML format.`
