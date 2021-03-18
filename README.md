# [<img src="https://ipinfo.io/static/ipinfo-small.svg" alt="IPinfo" width="24"/>](https://ipinfo.io/) IPinfo CLI

This is the official CLI for the [IPinfo.io](https://ipinfo.io) IP address API,
allowing you to:

- Look up IP details in bulk or one-by-one.
- Look up ASN details.
- Summarize the details of up to 1000 IPs at a time.
- Sign up for a IPinfo API token.

## Installation

TODO - give a curl/wget cmd for installing from github archives, unpacking that
and installing it to common locations for different operating systems.

## Quick Start

By default, invoking the CLI shows a help message:

```bash
$ ipinfo
Usage: ipinfo <cmd> [<opts>] [<args>]

Commands:
  <ip>        look up details for an IP address, e.g. 8.8.8.8.
  <asn>       look up details for an ASN, e.g. AS123 or as123.
  myip        get details for your IP.
  bulk        get details for multiple IPs in bulk.
  summarize   get summarized data for a group of IPs.
  prips       print IP list from CIDR or range.
  login       save an API token session.
  logout      delete your current API token session.
  version     show current version.

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --help, -h
      show help.

  Outputs:
    --field, -f
      lookup only a specific field in the output.
      field names correspond to JSON keys, e.g. 'hostname' or 'company.type'.

  Formats:
    --pretty, -p
      output pretty format.
    --json, -j
      output JSON format.
    --csv, -c
      output CSV format.
```

If you have a token, log in with it first. You can continue without a token,
but there will be limited data output and some features (like bulk lookups)
will not be available.

```bash
$ ipinfo login
```

You can quickly look up details of your own IP with `myip`:

```bash
$ ipinfo myip
```

Or of another IP by specifying it:

```bash
$ ipinfo 8.8.8.8
```

You can change the format of the output to JSON:

```bash
$ ipinfo 8.8.8.8 --json
```

And in case you only needed a single field:

```bash
$ ipinfo 8.8.8.8 -f hostname
```

And if you have the need to input IPs from `stdin`, just pipe it in (this will
require having a token!):

```bash
$ cat ip-list.txt | ipinfo --json
```

There are **many** more features available, so for full details, consult the
`-h` or `--help` message for each command. For example:

```bash
$ ipinfo 8.8.8.8 --help
```

## Other IPinfo Tools

There are official IPinfo client libraries available for many languages including PHP, Python, Go, Java, Ruby, and many popular frameworks such as Django, Rails and Laravel. There are also many third party libraries and integrations available for our API.

See [https://ipinfo.io/developers/libraries](https://ipinfo.io/developers/libraries) for more details.

## About IPinfo

Founded in 2013, IPinfo prides itself on being the most reliable, accurate, and in-depth source of IP address data available anywhere. We process terabytes of data to produce our custom IP geolocation, company, carrier and IP type data sets. Our API handles over 12 billion requests a month for 100,000 businesses and developers.

[![image](https://avatars3.githubusercontent.com/u/15721521?s=128&u=7bb7dde5c4991335fb234e68a30971944abc6bf3&v=4)](https://ipinfo.io/)
