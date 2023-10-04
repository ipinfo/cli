# 3.1.2

* `ipinfo tool prefix` introduced with some misc. prefix tools. Currently supports following:
    - `ipinfo tool prefix addr` returns the base IP address of a prefix.
    - `ipinfo tool prefix bits` returns the length of a prefix and reports `-1`` if invalid.
    - `ipinfo tool prefix masked` returns canonical form of a prefix, masking off non-high bits, and returns the zero if invalid.
    - `ipinfo tool prefix is_valid` reports whether a prefix is valid.

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
