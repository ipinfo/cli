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

The `ipinfo` CLI is available for download via multiple mechanisms.

### macOS

```bash
brew install ipinfo-cli
```

OR to install the latest `amd64` version without automatic updates:

```bash
curl -Ls https://github.com/ipinfo/cli/releases/download/ipinfo-3.1.2/macos.sh | sh
```

### Ubuntu PPA

_Note_: this installs our full suite of binaries and keeps them up-to-date.

```bash
sudo add-apt-repository ppa:ipinfo/ppa
sudo apt update
sudo apt install ipinfo
```

### Debian / Ubuntu (amd64)

_Note_: this is a one-time installation; updates are not automatic. Use the PPA
for automatic updates.

```bash
curl -Ls https://github.com/ipinfo/cli/releases/download/ipinfo-3.1.2/deb.sh | sh
```

OR

```bash
curl -LO https://github.com/ipinfo/cli/releases/download/ipinfo-3.1.2/ipinfo_3.1.2.deb
sudo dpkg -i ipinfo_3.1.2.deb
```

### FreeBSD

```bash
cd /usr/ports/net/ipinfo-cli && make install clean
```

### Arch linux

```bash
git clone https://aur.archlinux.org/ipinfo-cli.git
makepkg -si
```

### Windows Powershell

_Note_: run powershell as administrator before executing this command.

```bash
iwr -useb https://github.com/ipinfo/cli/releases/download/ipinfo-3.1.2/windows.ps1 | iex
```

### Scoop

```bash
scoop install ipinfo-cli
```

### Docker

```bash
docker run --rm -it ipinfo/ipinfo:3.1.2
```

To save the CLI's config, add `-v "/path_to_config:/root/.config/ipinfo"`. For
example, the following command saves the config to the `ipinfo` directory in
the current working directory.

```bash
docker run --rm -it -v "$PWD/ipinfo:/root/.config/ipinfo" ipinfo/ipinfo:3.1.2
```

### Using `go install`

Make sure that `$GOPATH/bin` is in your `$PATH`, because that's where this gets
installed:

```bash
go install github.com/ipinfo/cli/ipinfo@latest
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
solaris_amd64
windows_386
windows_amd64
windows_arm
```

After choosing a platform `PLAT` from above, run:

```bash
# for Windows, use ".zip" instead of ".tar.gz"
curl -LO https://github.com/ipinfo/cli/releases/download/ipinfo-3.1.2/ipinfo_3.1.2_${PLAT}.tar.gz
# OR
wget https://github.com/ipinfo/cli/releases/download/ipinfo-3.1.2/ipinfo_3.1.2_${PLAT}.tar.gz

tar -xvf ipinfo_3.1.2_${PLAT}.tar.gz
mv ipinfo_3.1.2_${PLAT} /usr/local/bin/ipinfo
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

## Additional CLIs

The `ipinfo` CLI has some subcommands like `grepip`, `prips`, `cidr2range`, `cidr2ip`,
`range2cidr`, `range2ip`, `splitcidr`, `randip` and `mmdb` which are also shipped as standalone binaries.

These binaries are available via all the **same installation methods** as
mentioned above for `ipinfo`, except you must change only the name to the name
of the subcommand, and choose the appropriate version.

Currently these subcommands are separately shipped:

| CLI        | Version                                                               |
| ---------- | --------------------------------------------------------------------- |
| grepip     | [1.2.2](https://github.com/ipinfo/cli/releases/tag/grepip-1.2.2)      |
| prips      | [1.0.0](https://github.com/ipinfo/cli/releases/tag/prips-1.0.0)       |
| cidr2range | [1.2.0](https://github.com/ipinfo/cli/releases/tag/cidr2range-1.2.0)  |
| cidr2ip    | [1.0.0](https://github.com/ipinfo/cli/releases/tag/cidr2ip-1.0.0)     |
| range2cidr | [1.3.0](https://github.com/ipinfo/cli/releases/tag/range2cidr-1.3.0)  |
| range2ip   | [1.0.0](https://github.com/ipinfo/cli/releases/tag/range2ip-1.0.0)    |
| randip     | [1.1.0](https://github.com/ipinfo/cli/releases/tag/randip-1.1.0)      |
| splitcidr  | [1.0.0](https://github.com/ipinfo/cli/releases/tag/splitcidr-1.0.0)   |
| mmdb       | [1.4.2](https://github.com/ipinfo/mmdbctl/releases/tag/mmdbctl-1.4.2) |

## Quick Start

This will help you quickly get started with the `ipinfo` CLI.

### Default Help Message

By default, invoking the CLI shows a help message:

```bash
ipinfo
```

![ipinfo](gif/default.gif)

### Login

If you have a token, log in with it first. You can continue without a token,
but there will be limited data output and some features (like bulk lookups)
will not be available. Get your token for free at
[https://ipinfo.io/signup](https://ipinfo.io/signup?ref=cli).

```bash
ipinfo init
```

### My IP

You can quickly look up details of your own IP with `myip`:

```bash
ipinfo myip
```

![ipinfo myip](gif/myip.gif)

### Any IP

You can see the details of any IP by specifying it:

```bash
ipinfo 8.8.8.8
```

![ipinfo myip](gif/ip8.8.8.8.gif)

### Piping

You can pipe IPs in and get their results in bulk (this requires a token):

```bash
cat ips.txt | ipinfo | less
```

![cat ips.txt | ipinfo](gif/cat.gif)

Here's the CSV version of that:

```bash
cat ips.txt | ipinfo -c | less
```

![cat ips.txt | ipinfo -c](gif/cat-csv.gif)

### Field Filter

In case you only needed a single field from a bunch of IPs:

```bash
cat ips.txt | ipinfo -f hostname
```

![cat ips.txt | ipinfo](gif/hostname.gif)

### Bulk

The above commands implicitly run the `bulk` subcommand on the input. You can
manually specify bulk and input IPs on the command line:

```bash
ipinfo bulk 1.1.1.0/30 8.8.8.0/30 9.9.9.0/30 | less
```

![ipinfo bulk](gif/bulk.gif)

### Summarize

IP details can be summarized similar to what's provided by
https://ipinfo.io/tools/summarize-ips:

```bash
cat lk-ips.txt | ipinfo summarize
```

![ipinfo summarize](gif/summarize.gif)

There are many more features available, so for full details, consult the `-h`
or `--help` message for each command. For example:

```bash
ipinfo 8.8.8.8 --help
```

## Auto-Completion

Auto-completion is supported for at least the following shells:

```
bash
zsh
fish
```

NOTE: it may work for other shells as well because the implementation is in
Golang and is not necessarily shell-specific.

### Installation

Installing auto-completions is as simple as running one command (works for
`bash`, `zsh` and `fish` shells):

```bash
ipinfo completion install
```

If you want to customize the installation process (e.g. in case the
auto-installation doesn't work as expected), you can request the actual
completion script for each shell:

```bash
# get bash completion script
ipinfo completion bash

# get zsh completion script
ipinfo completion zsh

# get fish completion script
ipinfo completion fish
```

### Shell not listed?

If your shell is not listed here, you can open an issue.

Note that as long as the `COMP_LINE` environment variable is provided to the
binary itself, it will output completion results. So if your shell provides a
way to pass `COMP_LINE` on auto-completion attempts to a binary, then have your
shell do that with the `ipinfo` binary itself (or any of our binaries).

## Data

The amount of data you get back per lookup depends upon how much data you have
enabled on your token via the https://ipinfo.io site.

If you have an account, see our
[plans](https://ipinfo.io/account/billing/upgrade) and
[addons](https://ipinfo.io/account/addons).

All examples in this document use a token with all data enabled.

## Color Output

### Disabling Color Output

All our CLIs respect either the `--nocolor` flag or the
[`NO_COLOR`](https://no-color.org/) environment variable to disable color
output.

### Color on Windows

To enable color support for the Windows command prompt, run the following to
enable [`Console Virtual Terminal Sequences`](https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences).

```cmd
REG ADD HKCU\CONSOLE /f /v VirtualTerminalLevel /t REG_DWORD /d 1
```

You can disable this by running the following:

```cmd
REG DELETE HKCU\CONSOLE /f /v VirtualTerminalLevel
```

## Other IPinfo Tools

There are official IPinfo client libraries available for many languages including PHP, Python, Go, Java, Ruby, and many popular frameworks such as Django, Rails and Laravel. There are also many third party libraries and integrations available for our API.

See [https://ipinfo.io/developers/libraries](https://ipinfo.io/developers/libraries) for more details.

## About IPinfo

Founded in 2013, IPinfo prides itself on being the most reliable, accurate, and in-depth source of IP address data available anywhere. We process terabytes of data to produce our custom IP geolocation, company, carrier, VPN detection, hosted domains, and IP type data sets. Our API handles over 40 billion requests a month for 100,000 businesses and developers.

[![image](https://avatars3.githubusercontent.com/u/15721521?s=128&u=7bb7dde5c4991335fb234e68a30971944abc6bf3&v=4)](https://ipinfo.io/)
