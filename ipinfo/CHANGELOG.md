# 3.3.1

- Added support for Windows ARM64.
- Windows users can also install via Winget, Chocolatey, and Scoop.
- Updated Ubuntu PPA source. New URL is https://ppa.ipinfo.net.
- Added multiple arch support to PPA. These architectures are; i386, amd64, armhf, and arm64.
- Fixed ipinfo tool aggregate not working properly for adjacent CIRDs.
- Added some basic IP tools such as:
  - is_loopback
  - is_multicast
  - is_unspecified
  - is_global_unicast
  - is_interface_local_multicast
  - is_link_local_multicast
  - is_link_local_unicast
- Fixed some issues related to convenience scripts.

## Pull Requests

- [#192](https://github.com/ipinfo/cli/pull/192)
- [#193](https://github.com/ipinfo/cli/pull/193)
- [#203](https://github.com/ipinfo/cli/pull/203)
- [#204](https://github.com/ipinfo/cli/pull/204)
- [#205](https://github.com/ipinfo/cli/pull/205)
- [#207](https://github.com/ipinfo/cli/pull/207)

# 3.3.0

- Support CIDRs & Ranges in `grepip`.
- New subcommand `matchip` which helps to grep for IP CIDRs/Ranges that
  overlap.
- New command `grepdomain` which is like `grepip` but for domains.
- Added `-6` flag to `myip`.
- Fixed a bug with the signup process; now works smoothly.

## Pull Requests

- [#183](https://github.com/ipinfo/cli/pull/183)
- [#184](https://github.com/ipinfo/cli/pull/184)
- [#189](https://github.com/ipinfo/cli/pull/189)
- [#188](https://github.com/ipinfo/cli/pull/188)
- [#185](https://github.com/ipinfo/cli/pull/185)
- [#190](https://github.com/ipinfo/cli/pull/190)

# 3.2.0

* `ipinfo tool prefix` introduced with some misc. prefix tools. Currently supports following subcommands:
    - `addr` returns the base IP address of a prefix.
    - `bits` returns the length of a prefix and reports `-1` if invalid.
    - `masked` returns canonical form of a prefix, masking off non-high bits, and returns the zero if invalid.
    - `is_valid` reports whether a prefix is valid.

## Pull Requests

* [#179](https://github.com/ipinfo/cli/pull/179)
* [#180](https://github.com/ipinfo/cli/pull/180)

# 3.1.1

* Fixed return errors in IP parsing generic funcs.

# 3.1.0

* IP calculator introduced in `calc` subcommand for doing arbitrary arithmetic
on IP addresses.
* Bulk ASN now supported via `asn bulk` subcommand.
* `ipinfo tool upper` introduced to get the upper IP of an IP range or CIDR.
* `ipinfo tool lower` introduced to get the lower IP of an IP range or CIDR.
* `ipinfo tool next` introduced to get the next IP of a given IP.
* `ipinfo tool prev` introduced to get the previous IP of a given IP.
* `ipinfo tool is_v4` introduced to check whether an IP is v4.
* `ipinfo tool is_v6` introduced to check whether an IP is v6.
* `ipinfo tool is_one_ip` introduced to check whether an IP range or CIDR
* contains only a single IP.
* `ipinfo tool is_valid` introduced to check whether an IP is a valid IP or
not.
* `ipinfo tool unmap` introduced to return an IP with any IPv4-mapped IPv6
address prefix removed.
* Fixed CLI token parameter not being recognized in `download` subcommand.
* Fixed CLI login not getting saved after init.
* Now performing a checksum comparison on database downloads.
* YAML output now has null values removed for cleaner output.

## Pull Requests

* [#157](https://github.com/ipinfo/cli/pull/157)
* [#158](https://github.com/ipinfo/cli/pull/158)
* [#160](https://github.com/ipinfo/cli/pull/160)
* [#154](https://github.com/ipinfo/cli/pull/154)
* [#161](https://github.com/ipinfo/cli/pull/161)
* [#159](https://github.com/ipinfo/cli/pull/159)
* [#162](https://github.com/ipinfo/cli/pull/162)
* [#155](https://github.com/ipinfo/cli/pull/155)
* [#165](https://github.com/ipinfo/cli/pull/165)
* [#164](https://github.com/ipinfo/cli/pull/164)
* [#172](https://github.com/ipinfo/cli/pull/172)
* [#169](https://github.com/ipinfo/cli/pull/169)
* [#170](https://github.com/ipinfo/cli/pull/170)
* [#173](https://github.com/ipinfo/cli/pull/173)
