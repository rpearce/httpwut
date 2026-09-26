# httpwut modernization design

Date: 2026-09-22
Status: implemented on branch modernize; amended after final review (F1-F5) and code review

This spec follows the same shape and conventions as the `ature` modernization
spec of the same date so the two CLIs stay consistent. Deliberate differences
are listed in section 15.

## 1. Goal

Bring `httpwut`, a small cobra-based CLI that looks up HTTP status codes, up to
current Go standards as of September 2026: modern toolchain and dependencies,
idiomatic and testable code, unit tests, lint, CI, and automated releases. Keep
it small and readable. It is a personal learning project, not a library.

Success criteria:

- Builds and tests on Go 1.27 with current Cobra.
- Status data is a pure, unit-tested package with no cobra dependency.
- Commands are constructed by a function and tested end to end through
  captured output.
- golangci-lint v2 passes with the repo config.
- GitHub Actions runs tests, lint, and govulncheck on push and pull request,
  and goreleaser v2 produces releases on `v*` tags.
- mise pins the Go toolchain and tools for anyone who clones the repo.

## 2. Non-goals

- No new user-facing features beyond `--version`, shell completion, and the
  data corrections listed in section 8.
- No Makefile or Taskfile, no CHANGELOG file, no Homebrew tap, no Docker, no
  SBOM or signing, no cross-OS CI matrix, no Renovate, no `tool` directives in
  go.mod.
- No embedded JSON/YAML data or code generation. The status table stays a Go
  map literal.
- No change to the default output format of `httpwut is <code>`.
- No `context` plumbing or signal handling; nothing here is cancellable.
- No new git tag or published release, and no pushing the branch. The next
  release is the owner's call.

## 3. Audit findings driving the design

Verified on 2026-09-22:

- go.mod declares `go 1.18`; cobra v1.8.1; `github.com/pkg/browser` at a 2024
  pseudo-version with no tags and a frozen dependency floor. Cobra v1.10.2 has
  no breaking changes affecting this code.
- All code lives in package `cmd` (non-main, which collides with the Go
  convention that `cmd/<name>` holds main packages) and in `main.go`.
- Output goes through `fmt.Println`, flags are read by name with ignored
  errors, args are checked manually, and the browser opener is called
  directly. None of this is unit-testable. There are no tests.
- Extra arguments are silently ignored: `httpwut is 200 500` prints `200 - OK`.
- `abc` and `999` produce the same "not found" error, so a typo and an
  unknown code are indistinguishable.
- Every error prints the full usage block to stderr.
- If `--cats` fails to open a browser, `--dogs` never runs.
- `--version` does not exist.
- No CI, no lint config, no Dependabot config. `.goreleaser.yaml` is
  `version: 2` but uses `archives.format`, deprecated since v2.6, and the
  locally installed goreleaser is v1.26.2 and refuses the file.
- The `go` shim on this machine errors because no Go version is pinned for the
  project (Go is managed by mise).
- Data bugs: 401 is titled "Not Authorized"; several descriptions have
  copy-paste damage; the 103 link points at the WHATWG HTML spec; eight IANA-
  registered codes are missing. README links to `http.cats`, which does not
  resolve.

## 4. Toolchain and dependencies

| Item | Decision |
|---|---|
| go.mod `go` directive | `go 1.26.0` plus `toolchain go1.27.1` |
| Local Go via mise | 1.27.1 |
| golangci-lint | 2.13.2, config format v2 |
| goreleaser | 2.18.2 |
| govulncheck | 1.8.0 |
| cobra | v1.10.2 |
| browser opener | `github.com/cli/browser` v1.3.0, replacing `github.com/pkg/browser` |

Rationale:

- `go 1.26.0` is the oldest supported Go release line, so `go install` works
  for anyone on a supported Go without a toolchain download. The dependencies
  only need Go 1.21, and `go vet`'s stdversion check stops code from using
  newer standard-library APIs. `toolchain go1.27.1` makes setup-go in CI and
  release install the same Go that `.mise.toml` pins; setup-go reads the
  `toolchain` line when present.
- Tools are pinned in a new `.mise.toml` at the repo root. This fixes the `go`
  shim error and the goreleaser v1/v2 mismatch. golangci-lint's own docs
  discourage `go tool` and `tool` directives for it, and mise is already the
  tool manager on this machine.
- `cli/browser` has the same `OpenURL` API as `pkg/browser`, a semver tag,
  maintained dependencies, WSL support, and the GitHub CLI as a downstream
  user.

`.mise.toml`:

```toml
[tools]
go = "1.27.1"
golangci-lint = "2.13.2"
goreleaser = "2.18.2"
"go:golang.org/x/vuln/cmd/govulncheck" = "1.8.0"

[tasks.test]
description = "Run the tests with the race detector"
run = "go test -race -cover -shuffle=on ./..."

[tasks.lint]
description = "Run golangci-lint"
run = "golangci-lint run"

[tasks.fmt]
description = "Format Go files with gofumpt and goimports"
run = "golangci-lint fmt"

[tasks.vuln]
description = "Check dependencies for known vulnerabilities"
run = "govulncheck ./..."

[tasks.snapshot]
description = "Build a local snapshot release into dist/"
run = "goreleaser release --snapshot --clean"
```

## 5. Layout

The `cmd/` package is deleted. `main.go` stays at the module root so
`go install github.com/rpearce/httpwut@latest` keeps working.

```
main.go                      package main: version var set by -ldflags,
                             prints "httpwut: <err>" to stderr, exits 1
main_test.go                 resolveVersion table test
internal/status/status.go    Status type, Lookup, All
internal/status/codes.go     the map literal
internal/status/status_test.go
internal/cli/root.go         NewRootCmd(version string) *cobra.Command
internal/cli/is.go           is command
internal/cli/cli_test.go
.mise.toml
.golangci.yml
.goreleaser.yaml             rewritten for v2
.github/workflows/ci.yml
.github/workflows/release.yml
.github/dependabot.yml
docs/superpowers/specs/      this document
```

Every package gets a package doc comment in one of its files.

## 6. Package `internal/status`

Pure data and lookup, no I/O, no dependencies.

```go
// Package status holds the table of HTTP status codes and their descriptions.
package status

// Status describes one HTTP status code.
type Status struct {
	Code        int
	Title       string
	Description string
	URL         string
}

// Lookup returns the status for code and whether it is known.
func Lookup(code int) (Status, bool)

// All returns every known status sorted by code.
func All() []Status
```

- The data is `var statuses = map[int]Status{...}` in `codes.go`. Entries in
  the literal set only `Title`, `Description`, and `URL`; the map key is the
  code. `Lookup` and `All` copy the key into `Code` before returning, so the
  code is never written twice.
- `All` builds a slice from the map and sorts it with `slices.SortFunc`. It is
  only called for completion, so no caching.

## 7. Package `internal/cli`

`NewRootCmd(version string) *cobra.Command` builds a fresh root on every
call. There is no package-level command state and no `init()` registration,
so tests can construct independent roots and run in parallel.

```go
// NewRootCmd builds the root command with the real browser opener.
func NewRootCmd(version string) *cobra.Command

func newRootCmd(version string, openURL func(url string) error) *cobra.Command
func newIsCmd(openURL func(url string) error) *cobra.Command
```

`openURL` defaults to `browser.OpenURL` from `cli/browser`. Tests pass a
stub.

Root command settings:

- `Use: "httpwut"`, `Short: "Look up HTTP status code information"`.
- No `Args` field. Cobra's default already rejects `httpwut 502` with
  `unknown command "502" for "httpwut"`. Setting `cobra.NoArgs` on a root
  that has no `Run` would make cobra print help instead, because it returns
  help for a non-runnable command before validating args.
- `Version: version`, which gives Cobra's built-in `--version` flag.
- `SilenceUsage: true` and `SilenceErrors: true`. `main` prints the error.
- `CompletionOptions.HiddenDefaultCmd: true` so `completion` still works but
  no longer appears in `-h` output.

`is` command:

- `Use: "is <code>"`, `Short: "Look up an HTTP status code"`.
- `Args: cobra.ExactArgs(1)`, so zero or two arguments is an error.
- Flags bound with `BoolVarP` into an unexported options struct:
  `-v/--verbose`, `-c/--cats`, `-d/--dogs`. Help text unchanged.
- `ValidArgsFunction` returns `cobra.CompletionWithDesc(code, title)` for
  every entry of `status.All()` with `ShellCompDirectiveNoFileComp`.
- `RunE`:
  1. `strconv.Atoi(args[0])`; on failure return
     `fmt.Errorf("invalid HTTP status code %q", args[0])`.
  2. `status.Lookup(code)`; on miss return
     `fmt.Errorf("unknown HTTP status code %d", code)`.
  3. Write `"%d - %s\n"` with code and title to `cmd.OutOrStdout()`.
  4. If verbose, write description then URL, one per line, same writer.
  5. If cats, call `openURL(catsURL + strconv.Itoa(code))`. If dogs, call
     `openURL(dogsURL + strconv.Itoa(code))`. Both are attempted even if the
     first fails. Each failure is wrapped with context, e.g.
     `fmt.Errorf("open http.cat: %w", err)`, and the result is
     `errors.Join(errs...)`.
- Base URLs are unexported constants.

## 8. Data corrections and additions

Titles:

| Code | Old | New |
|---|---|---|
| 401 | Not Authorized | Unauthorized |
| 418 | (Unused / I'm a teapot) | I'm a teapot |

Descriptions (rewritten text):

- 204: "The server has successfully fulfilled the request and there is no
  additional content to send in the response content."
- 305: "305 was defined in a previous version of HTTP/1.1 and is now
  deprecated."
- 306: "306 was defined in a previous version of HTTP/1.1, is no longer used,
  and the code is reserved."
- 405: "The method received in the request-line is known by the origin server
  but not supported by the target resource."
- 410: "Access to the target resource is no longer available at the origin
  server and this condition is likely to be permanent."
- 412: "One or more conditions given in the request header fields evaluated to
  false when tested on the server."
- 418: "The server refuses to brew coffee because it is, permanently, a
  teapot. This is a reference to the Hyper Text Coffee Pot Control Protocol
  from the 1998 and 2014 April Fools' RFCs. RFC 9110 lists 418 as unused."
- 510: "The policy for accessing the resource has not been met in the request.
  This code is obsolete: RFC 2774 has been moved to Historic status."

URLs:

- 103: `https://www.rfc-editor.org/rfc/rfc8297`
- 451: `https://www.rfc-editor.org/rfc/rfc7725#section-3` (consistency with
  the other rfc-editor links)

New entries, titles matching `net/http.StatusText`:

| Code | Title | Description | URL |
|---|---|---|---|
| 102 | Processing | An interim response used to inform the client that the server has accepted the complete request but has not yet completed it. Deprecated by RFC 4918. | https://www.rfc-editor.org/rfc/rfc2518#section-10.1 |
| 207 | Multi-Status | The message body that follows is by default an XML message and can contain a number of separate response codes, depending on how many sub-requests were made. | https://www.rfc-editor.org/rfc/rfc4918#section-11.1 |
| 208 | Already Reported | The members of a DAV binding have already been enumerated in a preceding part of the multistatus response, and are not being included again. | https://www.rfc-editor.org/rfc/rfc5842#section-7.1 |
| 226 | IM Used | The server has fulfilled a GET request for the resource, and the response is a representation of the result of one or more instance-manipulations applied to the current instance. | https://www.rfc-editor.org/rfc/rfc3229#section-10.4.1 |
| 423 | Locked | The source or destination resource of a method is locked. | https://www.rfc-editor.org/rfc/rfc4918#section-11.3 |
| 424 | Failed Dependency | The method could not be performed on the resource because the requested action depended on another action and that action failed. | https://www.rfc-editor.org/rfc/rfc4918#section-11.4 |
| 425 | Too Early | The server is unwilling to risk processing a request that might be replayed. | https://www.rfc-editor.org/rfc/rfc8470#section-5.2 |
| 428 | Precondition Required | The origin server requires the request to be conditional. | https://www.rfc-editor.org/rfc/rfc6585#section-3 |

Everything else in the table is unchanged.

## 9. `main.go`

```go
var version string // set by goreleaser via -X main.version=...

func main() {
	if err := cli.NewRootCmd(buildVersion()).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "httpwut:", err)
		os.Exit(1)
	}
}

func buildVersion() string {
	info, _ := debug.ReadBuildInfo()
	return resolveVersion(version, info)
}

func resolveVersion(ldflag string, info *debug.BuildInfo) string
```

Errors are passed through `errorLine`, which prefixes `httpwut: ` to each
error of an `errors.Join` result and trims
trailing newlines, because cobra's unknown-command suggestions end in one.

`resolveVersion` order: `ldflag` if non-empty; else `info.Main.Version` if
`info` is non-nil and the value is non-empty and not `(devel)`; else `"dev"`.
Go 1.24+ stamps `Main.Version` from the VCS tag for `go build`, and
`go install pkg@vX` has always set it, so the ldflag is only needed for
goreleaser. Keeping the decision in a pure function makes it testable without
controlling the test binary's build info.

## 10. Tests

Tests are written before the implementation they cover. All tests use the
standard `testing` package with `t.Parallel()` on the test and each subtest.
No mocking frameworks, no testscript, no coverage gate.

### `internal/status/status_test.go`

Tests live in package `status`, not `status_test`, so they can see the map.

- `Lookup` table test: known code returns the entry with matching `Code`;
  unknown code (e.g. 999) returns `false`.
- Sweep over `All()`: every entry has non-empty `Title`, non-empty
  `Description`, and `URL` starting with `https://`. `Lookup(s.Code)` for each
  entry returns an equal `Status`, which checks that `Code` is filled from the
  map key. `len(All()) == len(statuses)`.
- `All()` is sorted ascending and has no duplicate codes.
- Title parity with `net/http`: for each entry, if `http.StatusText(code)` is
  non-empty and the code is not in the allowlist `{413, 414, 416, 422}`,
  assert `Title == http.StatusText(code)`. The allowlist exists because
  `net/http` keeps RFC 7231 wording for those four while this table uses RFC
  9110 wording. This test would have caught the 401 bug.
- Explicit presence test for the eight new codes.

### `internal/cli/cli_test.go`

End to end through `newRootCmd(version, stub)`, using `SetArgs`, `SetOut`,
`SetErr`, and `Execute`, asserting on exact strings.

- `is 502`: stdout is exactly `"502 - Bad Gateway\n"`, no error, stub not
  called.
- `is 502 --verbose`: stdout is the code line, description line, URL line.
- `is abc`: error message contains `invalid HTTP status code "abc"`.
- `is 999`: error message is `unknown HTTP status code 999`.
- `is` with no args and `is 200 500`: error mentions `accepts 1 arg(s)`.
- `is 404 --cats --dogs`: stub records exactly
  `["https://http.cat/404", "https://httpstatusdogs.com/404"]` in order.
- `is 404 --cats --dogs` with a stub that fails for the cat URL: both URLs
  still recorded, error returned wraps the stub error, stdout still has the
  code line.
- `--version` prints `httpwut version <v>`.
- `httpwut 502` (no subcommand) returns an error.
- Nothing is written to the error writer on the happy path, and nothing is
  written to the output writer on the error path, except the browser-failure
  case above where the lookup line is printed before the opener runs.

### `main_test.go`

- `resolveVersion` table test: ldflag set wins; nil info gives `"dev"`;
  `Main.Version` of `""` or `(devel)` gives `"dev"`; `Main.Version` of
  `v0.2.0` is returned as is.

## 11. Tooling files

Versions verified on 2026-09-22.

### `.golangci.yml`

```yaml
version: "2"

run:
  timeout: 3m

linters:
  default: standard
  enable:
    - copyloopvar
    - errorlint
    - exhaustive
    - gocritic
    - gosec
    - intrange
    - misspell
    - modernize
    - nilerr
    - paralleltest
    - perfsprint
    - predeclared
    - revive
    - thelper
    - tparallel
    - unconvert
    - unparam
    - usestdlibvars
  exclusions:
    presets:
      - std-error-handling

formatters:
  enable:
    - gofumpt
    - goimports
  settings:
    gofumpt:
      extra:
        group-params: true
    goimports:
      local-prefixes:
        - github.com/rpearce/httpwut
```

`standard` is errcheck, govet, ineffassign, staticcheck, unused. This is the
`ature` set plus `copyloopvar` and `modernize`. revive's default rules stay
enabled, including package and exported-symbol doc comments, which the code
will satisfy. `golangci-lint run` and `golangci-lint fmt --diff` must be
clean locally and in CI.

### `.goreleaser.yaml`

```yaml
version: 2

builds:
  - main: .
    binary: httpwut
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    mod_timestamp: "{{ .CommitTimestamp }}"
    flags:
      - -trimpath
    ldflags:
      - -s -w -X main.version=v{{ .Version }}

archives:
  - formats:
      - tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        formats:
          - zip

checksum:
  name_template: checksums.txt

snapshot:
  version_template: "{{ incpatch .Version }}-next"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
```

Changes from today: the `go mod tidy` and `go generate` hooks, example
comments, and modelines are removed (CI enforces tidiness, and a hook would
rewrite go.mod after goreleaser checked the tree is clean); deprecated
`format` keys become `formats`; 386 builds are dropped; the build gets
`-trimpath`, a reproducible `mod_timestamp`, and a version ldflag that
targets a symbol that exists. Archive names change from
`httpwut_Darwin_arm64.tar.gz` to `httpwut_0.2.0_darwin_arm64.tar.gz`,
matching `ature` and the goreleaser default. Nothing depends on the old
names.

### `.github/workflows/ci.yml`

On push to `main` and `pull_request`. `permissions: contents: read`. Runs on
`ubuntu-latest` only.

- `test` job: `actions/checkout@v7`, `actions/setup-go@v7` with
  `go-version-file: go.mod`, then
  `go mod tidy -diff`, `go vet ./...`,
  `go test -race -cover -shuffle=on ./...`.
- `lint` job: checkout, setup-go, `golangci/golangci-lint-action@v9` with
  `version: v2.13.2`.
- `vuln` job: checkout, setup-go with `go-version-file: go.mod`, then
  `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...`. This matches the
  mise pin; `golang/govulncheck-action` always installs the latest stable Go
  and govulncheck regardless of `go-version-file`.

### `.github/workflows/release.yml`

On push of tags matching `v*`. `permissions: contents: write`.
`actions/checkout@v7` with `fetch-depth: 0`, `actions/setup-go@v7` with
`go-version-file: go.mod`, then `go mod tidy -diff`, `go vet ./...`, and
`go test -race ./...` so a
tag on an untested commit cannot publish, then
`goreleaser/goreleaser-action@v7` with `distribution: goreleaser`,
`version: v2.18.2` (the mise pin), `args: release --clean`, and
`GITHUB_TOKEN` from `secrets.GITHUB_TOKEN`.

### `.github/dependabot.yml`

Weekly `gomod` and `github-actions` updates, both `directory: /`.

### `README.md`

- `http.cats` becomes `http.cat`; "at the time same" becomes "at the same
  time".
- Refresh the `httpwut is -h` block and the usage examples after the changes.
- Document `--version`, `httpwut completion <shell>`, and that release
  archives are attached to GitHub releases.
- Add a CI status badge.
- Keep the "tool for me to learn Go" note; it is honest.

`.gitignore` is already correct. `LICENSE` is untouched.

The committed config files also carry short explanatory comments that this
spec does not reproduce.

## 12. Behavior compatibility

Unchanged:

- `httpwut is 502` prints `502 - Bad Gateway`.
- `--verbose`, `--cats`, `--dogs` semantics and short flags.
- Exit code 1 on any error.

Changed, all intentional:

- Errors print as `httpwut: <message>` on stderr, without the usage block.
  `-h` still shows usage.
- Extra positional args to `is` are an error instead of being ignored.
- Non-numeric input gets `invalid HTTP status code "abc"`; unknown numeric
  input gets `unknown HTTP status code 999`. The old message was
  `HTTP status code not found: "999"`.
- New `--version` flag. The `completion` subcommand still works but is hidden
  from `-h`.
- `--cats` and `--dogs` both run even if one fails.
- Data corrections and eight new codes per section 8.
- Release archive names and the set of targets change per section 11.

## 13. Verification

Before the work is called done, all of the following pass locally under
mise:

```
go build ./...
go vet ./...
go test -race -cover -shuffle=on ./...
golangci-lint run
golangci-lint fmt --diff
govulncheck ./...
goreleaser check
goreleaser release --snapshot --clean   # dist/.../httpwut --version -> httpwut version v0.1.2-next
go run . --version
go run . is 502                 # 502 - Bad Gateway
go run . is 502 -v              # three lines
go run . is abc                 # httpwut: invalid HTTP status code "abc", exit 1
go run . is 999                 # httpwut: unknown HTTP status code 999, exit 1
go run . is 200 500             # error, exit 1, no usage block
go run . 502                    # error, exit 1
```

CI must be green on the pull request, and README examples must match actual
program output.

## 14. Out of scope

- Any new command or CLI shape.
- Rewriting the data source to IANA CSV, JSON, or YAML.
- Creating a new git tag or publishing a release.
- Pushing the branch.

## 15. Deliberate differences from the `ature` spec

Kept on purpose; each is a small improvement `ature` could adopt too.

- A `vuln` CI job running the pinned govulncheck, and a matching mise
  task. `ature` has no vulnerability check.
- `go mod tidy -diff` instead of `go mod tidy && git diff --exit-code`. Same
  result, one command, no working-tree mutation.
- A `toolchain` line in go.mod so CI and releases use the mise-pinned Go.
- Version resolution is a pure `resolveVersion` function with a table test in
  `main_test.go`. `ature` leaves `buildVersion` untested.
- Two extra linters, `copyloopvar` and `modernize`.
- The `openURL` injection point, which `ature` does not need.
