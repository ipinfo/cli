# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository overview

This is the official IPinfo CLI suite, written in Go. The repo produces one
"main" binary (`ipinfo`) plus several **standalone** binaries that are also
exposed as subcommands of `ipinfo`: `grepip`, `grepdomain`, `matchip`,
`prips`, `cidr2range`, `cidr2ip`, `range2cidr`, `range2ip`, `splitcidr`,
`randip`, and `mmdb`.

Each standalone CLI lives in its own top-level directory (e.g. `grepip/`,
`prips/`) containing only a thin `main.go` plus packaging scripts. The actual
implementation of every command — including the `ipinfo` subcommands that
share behavior with the standalone CLIs — lives under `lib/`. The `ipinfo/`
directory contains the `ipinfo`-specific commands (API calls, bulk lookups,
caching, config, mmdb tooling, etc.) and wires the shared `lib/` commands in
as subcommands.

## Common commands

```bash
# Build a single CLI into ./build/<name>
./scripts/build.sh ipinfo
./scripts/build.sh grepip
# (works for any top-level CLI dir: ipinfo, grepip, grepdomain, matchip,
#  prips, cidr2range, cidr2ip, range2cidr, range2ip, splitcidr, randip)

# Build a CLI for all supported OS/arch combos (used by release tooling)
./scripts/build-all-platforms.sh ipinfo <version>

# Format / lint everything in the project
./scripts/fmt.sh
./scripts/lint.sh   # requires `golint` on PATH

# Tests (only lib/iputil and lib/complete have tests)
go test ./lib/...
go test ./lib/iputil -run TestSubnet   # single test by name
```

Go module is vendored (`vendor/`); `go build` / `go test` use the vendored
deps automatically. Minimum Go version is in `go.mod`.

## Architecture

### Three layers

1. **`lib/iputil/`** — pure IP/CIDR/range primitives (parsing, conversion,
   math, IPv4/IPv6 helpers, u128, regexes for IP detection, file/stdin
   helpers). No CLI concerns. This is the foundation everything else builds
   on.
2. **`lib/`** — command implementations shared between the standalone
   binaries and `ipinfo`'s subcommands. Each `lib/cmd_*.go` file exposes a
   `CmdXOpts` struct, a `printHelp` function, and a `CmdX(...)` entrypoint
   that parses flags via `pflag` and runs the command. Standalone CLIs in
   `<name>/main.go` do little more than parse top-level flags and call into
   `lib`.
3. **`ipinfo/`** — the umbrella `ipinfo` binary. `ipinfo/main.go` dispatches
   on `os.Args[1]`: bare IPs/ASNs/domains short-circuit to lookup commands;
   otherwise it routes to a `cmd_*.go` file in this directory. The
   subcommands that mirror standalone tools (e.g. `cmd_grepip.go`,
   `cmd_prips.go`) are thin wrappers that delegate to `lib`.

### `ipinfo/`-specific concerns

- **API client** (`ipinfo_client.go`) wraps `github.com/ipinfo/go/v2`.
- **Cache** (`cache.go`) uses `bbolt` on disk to memoize API responses;
  the `--no-cache` flag and `cmd_cache.go` manage it.
- **Config / auth** (`config.go`, `cmd_config.go`, `init.go`,
  `cmd_init.go`, `cmd_logout.go`) persists tokens and settings under
  `~/.config/ipinfo`.
- **MMDB tooling** (`cmd_mmdb*.go`) wraps `mmdbctl` for read/import/export/
  diff/verify/metadata of MaxMind-format databases.
- **Shell completion**: each binary has a `completions.go` that registers
  via `lib/complete`. The `lib/complete/` tree is a vendored fork of
  `posener/complete` with project-local extensions.

### Adding a new command

- If it's pure IP math/string logic, add it under `lib/` as `cmd_foo.go`
  following the existing `CmdFoo` / `CmdFooOpts` / `printHelpFoo` pattern,
  then expose it from `ipinfo/cmd_foo.go` (and optionally a new top-level
  standalone dir with its own `main.go` + packaging scripts).
- Update `scripts/lint.sh` and `scripts/fmt.sh` if you add a new top-level
  CLI dir — they list each one explicitly.
- Run `./scripts/fmt.sh` before committing.

### Releases

`scripts/release.sh <cli> <version>` builds for all platforms via
`build-archive-all.sh` and uploads to GitHub. Per-CLI `release.sh`,
`deb.sh`, `docker.sh`, `macos.sh`, `windows.ps1` files are packaging
helpers consumed by the release pipeline; they are not for local dev.
