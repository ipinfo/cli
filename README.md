# [<img src="https://ipinfo.io/static/ipinfo-small.svg" alt="IPinfo" width="24"/>](https://ipinfo.io/) IPinfo CLI

This is the official CLI for the [IPinfo.io](https://ipinfo.io) IP address API,
allowing you to:

- Look up IP details in bulk or one-by-one.
- Look up ASN details.
- Summarize the details of up to 1000 IPs at a time.
- Open a map of IP locations for any set of IPs.
- Filter IPv4 & IPv6 addresses from any input.
- Print out IP lists for any CIDR or IP range.
- And more!

## Installation

All CLI tools (e.g. `ipinfo`, `grepip`) are available for download via
multiple mechanisms.

### macOS

```bash
brew install ipinfo-cli
```

OR to install the latest `amd64` version without automatic updates:

```bash
curl -Ls https://github.com/ipinfo/cli/releases/download/ipinfo-1.1.5/macos.sh | sh
```

### Debian / Ubuntu (amd64)

```bash
curl -Ls https://github.com/ipinfo/cli/releases/download/ipinfo-1.1.5/deb.sh | sh
```

OR

```bash
curl -LO https://github.com/ipinfo/cli/releases/download/ipinfo-1.1.5/ipinfo_1.1.5.deb
sudo dpkg -i ipinfo_1.1.5.deb
```

### Using `go get`

Make sure that `$GOPATH/bin` is in your `$PATH`, because that's where this gets
installed:

```bash
go get github.com/ipinfo/cli/ipinfo
```

### Using `curl`/`wget`

The pre-built binaries for all platforms are available on GitHub via artifacts
in releases. You need to simply download, unpack and move them to your shell's
binary search path.

The following OS & arch combinations are supported (if you use one not listed
on here, please open an issue):

```
darwin_amd64
darwin_arm64
dragonfly_amd64
freebsd_386
freebsd_amd64
freebsd_arm
freebsd_arm64
linux_386
linux_amd64
linux_arm
linux_arm64
netbsd_386
netbsd_amd64
netbsd_arm
netbsd_arm64
openbsd_386
openbsd_amd64
openbsd_arm
openbsd_arm64
plan9_386
plan9_amd64
plan9_arm
solaris_amd64
windows_386
windows_amd64
windows_arm
```

After choosing a platform `PLAT` from above, run:

```bash
# for Windows, use ".zip" instead of ".tar.gz"
curl -LO https://github.com/ipinfo/cli/releases/download/ipinfo-1.1.5/ipinfo_1.1.5_${PLAT}.tar.gz
# OR
wget https://github.com/ipinfo/cli/releases/download/ipinfo-1.1.5/ipinfo_1.1.5_${PLAT}.tar.gz

tar -xvf ipinfo_1.1.5_${PLAT}.tar.gz
mv ipinfo_1.1.5_${PLAT} /usr/local/bin/ipinfo
```

### Using `git`

Installing from source requires at least the Golang version specified in
`go.mod`. You can install the Golang toolchain from
[the official site](https://golang.org/doc/install).

Once the correct Golang version is installed, simply clone the repository and
install the binary:

```bash
git clone https://github.com/ipinfo/cli ipinfo-cli
cd ipinfo-cli
go install ./ipinfo/
$GOPATH/bin/ipinfo
```

You can add `$GOPATH/bin` to your `$PATH` to access `ipinfo` directly from
anywhere.

Alternatively, you can do the following to output the binary somewhere
specific:

```bash
git clone https://github.com/ipinfo/cli ipinfo-cli
cd ipinfo-cli
go build -o <path> ./ipinfo/
```

Replace `<path>` with the required location.

## Quick Start

### Default Help Message

By default, invoking the CLI shows a help message:

![ipinfo](gif/default.gif)

### Login

If you have a token, log in with it first. You can continue without a token,
but there will be limited data output and some features (like bulk lookups)
will not be available. Get your token for free at
[https://ipinfo.io/signup](https://ipinfo.io/signup?ref=cli).

```bash
ipinfo login
```

### My IP

You can quickly look up details of your own IP with `myip`:

![ipinfo myip](gif/myip.gif)

### Any IP

You can see the details of any IP by specifying it:

![ipinfo myip](gif/ip8.8.8.8.gif)

### Piping

You can pipe IPs in and get their results in bulk (this requires a token):

![cat ips.txt | ipinfo](gif/cat.gif)

Here's the CSV version of that:

![cat ips.txt | ipinfo -c](gif/cat-csv.gif)

### Field Filter

In case you only needed a single field from a bunch of IPs:

![cat ips.txt | ipinfo](gif/hostname.gif)

### Bulk

The above commands implicitly run the `bulk` subcommand on the input. You can
manually specify bulk and input IPs on the command line:

![ipinfo bulk](gif/bulk.gif)

### Summarize

IP details can be summarized similar to what's provided by
https://ipinfo.io/summarize-ips:

![ipinfo summarize](gif/summarize.gif)

There are many more features available, so for full details, consult the `-h`
or `--help` message for each command. For example:

```bash
ipinfo 8.8.8.8 --help
```

## Auto-Completion

Installing auto-completions is as simple as running one command (works for
`bash`, `zsh` and `fish` shells):

```bash
ipinfo completion
```

## Data

The amount of data you get back per lookup depends upon how much data you have
enabled on your token via the https://ipinfo.io site.

If you have an account, see our
[plans](https://ipinfo.io/account/billing/upgrade) and
[addons](https://ipinfo.io/account/addons).

All examples in this document use a token with all data enabled.

## Disabling Color Output

All our CLIs respect either the `--nocolor` flag or the
[`NO_COLOR`](https://no-color.org/)  environment variable to disable color
output.

## Other IPinfo Tools

There are official IPinfo client libraries available for many languages including PHP, Python, Go, Java, Ruby, and many popular frameworks such as Django, Rails and Laravel. There are also many third party libraries and integrations available for our API.

See [https://ipinfo.io/developers/libraries](https://ipinfo.io/developers/libraries) for more details.

## About IPinfo

Founded in 2013, IPinfo prides itself on being the most reliable, accurate, and in-depth source of IP address data available anywhere. We process terabytes of data to produce our custom IP geolocation, company, carrier, VPN detection, hosted domains, and IP type data sets. Our API handles over 40 billion requests a month for 100,000 businesses and developers.

[![image](https://avatars3.githubusercontent.com/u/15721521?s=128&u=7bb7dde5c4991335fb234e68a30971944abc6bf3&v=4)](https://ipinfo.io/)
