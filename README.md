# httpwut

[![CI](https://github.com/rpearce/httpwut/actions/workflows/ci.yml/badge.svg)](https://github.com/rpearce/httpwut/actions/workflows/ci.yml)

CLI tool to look up HTTP status code information.

_This is a tool for me to learn Go, so you probably shouldn't use this._

## Installation

```
go install github.com/rpearce/httpwut@latest
```

Prebuilt archives for Linux, macOS, and Windows are attached to each
[GitHub release](https://github.com/rpearce/httpwut/releases).

## Usage

```
λ httpwut is -h
Look up an HTTP status code

Usage:
  httpwut is <code> [flags]

Flags:
  -c, --cats      Open HTTP Status Cats webpage for status
  -d, --dogs      Open HTTP Status Dogs webpage for status
  -h, --help      help for is
  -v, --verbose   Print status code description and URL

λ httpwut is 502
502 - Bad Gateway

λ httpwut is 502 --verbose
502 - Bad Gateway
The server, while acting as a gateway or proxy, received an invalid response from an inbound server it accessed while attempting to fulfill the request.
https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway
```

Passing `--cats` and/or `--dogs` (you can do both at the same time) will open
the status code on https://http.cat and https://httpstatusdogs.com,
respectively.

### Version

```
λ httpwut --version
httpwut version v0.2.0
```

### Shell completion

`httpwut` can complete status codes for bash, zsh, fish, and PowerShell.
For zsh, for example:

```
λ httpwut completion zsh > "${fpath[1]}/_httpwut"
```

Run `httpwut completion --help` for the other shells.

## Development

Tools are pinned in `.mise.toml`; run `mise install` once. Then:

```
mise run test       # go test -race -cover -shuffle=on ./...
mise run lint       # golangci-lint run
mise run fmt        # golangci-lint fmt
mise run vuln       # govulncheck ./...
mise run snapshot   # goreleaser release --snapshot --clean
```

Releases are cut by pushing a `v*` tag; GitHub Actions runs goreleaser.
