# httpwut Modernization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bring the `httpwut` CLI to September 2026 Go standards: pinned toolchain, current deps, a testable package layout with unit tests, lint, CI, and goreleaser v2 releases, with the status-code data corrected and expanded.

**Architecture:** The cobra-cli `cmd/` package is replaced by two internal packages: `internal/status` (pure data plus `Lookup`/`All`, no cobra) and `internal/cli` (commands built by constructors, output written to the cobra writer, browser opener injected). `main.go` resolves the version, runs the root command, prints `httpwut: <err>` to stderr, and exits 1. Tooling is pinned in `.mise.toml`; CI and release run on GitHub Actions.

**Tech Stack:** Go 1.27, github.com/spf13/cobra v1.10.2, github.com/cli/browser v1.3.0, golangci-lint 2.13.2, goreleaser 2.18.2, govulncheck 1.8.0, mise, GitHub Actions.

**Spec:** `docs/superpowers/specs/2026-09-22-modernization-design.md`

## Global Constraints

- go.mod: `go 1.27.0`, no `toolchain` line. Local Go via mise is 1.27.1.
- Dependencies: exactly `github.com/spf13/cobra v1.10.2` and `github.com/cli/browser v1.3.0` as direct requires. No other new modules.
- Tool pins: golangci-lint `2.13.2`, goreleaser `2.18.2`, govulncheck `1.8.0`.
- Output format of `httpwut is <code>` is exactly `<code> - <Title>` followed by a newline. Verbose adds the description line and the URL line.
- Errors reach the user as `httpwut: <message>` on stderr, exit code 1, never followed by the usage block.
- Error wording: `invalid HTTP status code %q` for non-numeric input, `unknown HTTP status code %d` for numeric input not in the table.
- Every package has a package doc comment in exactly one of its files (`main.go`, `internal/status/status.go`, `internal/cli/root.go`). Every exported symbol has a doc comment.
- Every test function and subtest calls `t.Parallel()` first.
- Cobra root: `SilenceUsage: true`, `SilenceErrors: true`, `CompletionOptions.HiddenDefaultCmd: true`, and no `Args` field. Cobra's default arg handling already rejects `httpwut 502` with `unknown command "502" for "httpwut"`; setting `Args: cobra.NoArgs` on a root that has no `Run` would make cobra print help instead, because it returns help for a non-runnable command before validating args.
- All commands in this plan run from the repo root `/Users/bobert/projects/httpwut` on branch `modernize`. After Task 1 the `go`, `golangci-lint`, `goreleaser`, and `govulncheck` commands resolve through mise shims because `.mise.toml` exists.
- Every commit message ends with the line `Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>`.
- Do not push, tag, or create a release.

## Review Focus

Inputs the spec implies but does not enumerate. Each has a pinned test in the task named.

1. A negative code given after `--` (`httpwut is -- -1`) must produce `unknown HTTP status code -1`, not a flag-parse message. Pinned in Task 3, `TestIs` case "negative".
2. A number too large for an int (`httpwut is 99999999999999999999`) must produce `invalid HTTP status code "99999999999999999999"`, not a panic or an overflow message. Pinned in Task 3, `TestIs` case "overflow".
3. Running `httpwut` with no arguments must print help to stdout and exit 0, and `httpwut is -h` likewise. Pinned in Task 3, `TestRoot` subtests "bare prints help" and "is help".
4. Shell completion must filter by the typed prefix (`__complete is 40` returns only 400 through 409). Pinned in Task 3, `TestCompletion`.
5. A code with leading zeros (`httpwut is 0404`) is accepted as 404. This is deliberate leniency from `strconv.Atoi`, pinned so it cannot change by accident. Pinned in Task 3, `TestIs` case "leading zeros accepted".

---

### Task 1: Pin the toolchain and upgrade dependencies

**Files:**
- Create: `.mise.toml`
- Modify: `go.mod`, `go.sum`
- Modify: `cmd/is.go:6` (import path only; the file is deleted in Task 4)

**Interfaces:**
- Consumes: nothing.
- Produces: a working `go` shim in this directory; `github.com/cli/browser` in go.mod with `browser.OpenURL(url string) error`; cobra v1.10.2.

- [ ] **Step 1: Create `.mise.toml`**

```toml
[tools]
go = "1.27.1"
golangci-lint = "2.13.2"
goreleaser = "2.18.2"
"go:golang.org/x/vuln/cmd/govulncheck" = "1.8.0"

[tasks.test]
run = "go test -race -cover -shuffle=on ./..."

[tasks.lint]
run = "golangci-lint run"

[tasks.fmt]
run = "golangci-lint fmt"

[tasks.vuln]
run = "govulncheck ./..."

[tasks.snapshot]
run = "goreleaser release --snapshot --clean"
```

- [ ] **Step 2: Trust the config and install the tools**

Run: `mise trust && mise install`
Expected: output ends without errors; goreleaser 2.18.2 and govulncheck 1.8.0 download (go 1.27.1 and golangci-lint 2.13.2 may already be present).

- [ ] **Step 3: Verify the shims resolve in this directory**

Run: `go version && golangci-lint --version && goreleaser --version | head -3 && govulncheck -version`
Expected: `go version go1.27.1 darwin/arm64`; golangci-lint reports `2.13.2`; goreleaser reports `2.18.2`; govulncheck reports `v1.8.0`.

- [ ] **Step 4: Bump the go directive and swap the browser import**

Run:

```bash
go mod edit -go=1.27.0
sed -i '' 's#"github.com/pkg/browser"#"github.com/cli/browser"#' cmd/is.go
go get github.com/spf13/cobra@v1.10.2 github.com/cli/browser@v1.3.0
go mod tidy
```

Expected: `go get` prints upgrade lines for cobra, pflag, and x/sys and an add line for cli/browser. `go mod tidy` removes `github.com/pkg/browser`.

- [ ] **Step 5: Verify go.mod**

Run: `cat go.mod`
Expected content (indirect versions may be newer; there must be no `toolchain` line and no `github.com/pkg/browser`):

```
module github.com/rpearce/httpwut

go 1.27.0

require (
	github.com/cli/browser v1.3.0
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10 // indirect
	golang.org/x/sys v0.46.0 // indirect
)
```

Go 1.27's tidy may keep two `require` blocks instead of one. Either shape is fine.

- [ ] **Step 6: Build and smoke-run the old code on the new toolchain**

Run: `go build ./... && go vet ./... && go run . is 502`
Expected: no build or vet output, then `502 - Bad Gateway`.

- [ ] **Step 7: Commit**

```bash
git add .mise.toml go.mod go.sum cmd/is.go
git commit -m "build: pin toolchain with mise, upgrade to Go 1.27 and current deps

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 2: `internal/status` package with corrected and expanded data

**Files:**
- Create: `internal/status/status.go`
- Create: `internal/status/codes.go`
- Test: `internal/status/status_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces:
  - `type Status struct { Code int; Title, Description, URL string }`
  - `func Lookup(code int) (Status, bool)`
  - `func All() []Status` sorted ascending by `Code`.

- [ ] **Step 1: Write the failing tests**

Create `internal/status/status_test.go`:

```go
package status

import (
	"net/http"
	"slices"
	"strings"
	"testing"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		code      int
		wantTitle string
		wantOK    bool
	}{
		{name: "known", code: 502, wantTitle: "Bad Gateway", wantOK: true},
		{name: "lowest", code: 100, wantTitle: "Continue", wantOK: true},
		{name: "unknown", code: 999, wantOK: false},
		{name: "negative", code: -1, wantOK: false},
		{name: "zero", code: 0, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := Lookup(tt.code)
			if ok != tt.wantOK {
				t.Fatalf("Lookup(%d) ok = %v, want %v", tt.code, ok, tt.wantOK)
			}
			if !ok {
				if got != (Status{}) {
					t.Fatalf("Lookup(%d) = %+v, want zero Status", tt.code, got)
				}
				return
			}
			if got.Code != tt.code {
				t.Errorf("Lookup(%d).Code = %d, want %d", tt.code, got.Code, tt.code)
			}
			if got.Title != tt.wantTitle {
				t.Errorf("Lookup(%d).Title = %q, want %q", tt.code, got.Title, tt.wantTitle)
			}
		})
	}
}

func TestAllEntriesAreComplete(t *testing.T) {
	t.Parallel()

	all := All()
	if len(all) != len(statuses) {
		t.Fatalf("All() has %d entries, statuses has %d", len(all), len(statuses))
	}
	for _, s := range all {
		if s.Title == "" {
			t.Errorf("%d: empty Title", s.Code)
		}
		if s.Description == "" {
			t.Errorf("%d: empty Description", s.Code)
		}
		if !strings.HasPrefix(s.URL, "https://") {
			t.Errorf("%d: URL %q does not start with https://", s.Code, s.URL)
		}
		got, ok := Lookup(s.Code)
		if !ok || got != s {
			t.Errorf("Lookup(%d) = %+v, %v; want %+v, true", s.Code, got, ok, s)
		}
	}
}

func TestAllIsSortedAndUnique(t *testing.T) {
	t.Parallel()

	all := All()
	for i := 1; i < len(all); i++ {
		if all[i-1].Code >= all[i].Code {
			t.Fatalf("All() not strictly ascending at index %d: %d then %d", i, all[i-1].Code, all[i].Code)
		}
	}
}

// rfc9110Wording lists codes whose RFC 9110 title differs from the RFC 7231
// wording that net/http keeps for compatibility.
var rfc9110Wording = []int{413, 414, 416, 422}

func TestTitlesMatchNetHTTP(t *testing.T) {
	t.Parallel()

	for _, s := range All() {
		want := http.StatusText(s.Code)
		if want == "" || slices.Contains(rfc9110Wording, s.Code) {
			continue
		}
		if s.Title != want {
			t.Errorf("%d: Title = %q, want %q (net/http)", s.Code, s.Title, want)
		}
	}
}

func TestNewCodesPresent(t *testing.T) {
	t.Parallel()

	for _, code := range []int{102, 207, 208, 226, 423, 424, 425, 428} {
		if _, ok := Lookup(code); !ok {
			t.Errorf("Lookup(%d) missing", code)
		}
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./internal/status/`
Expected: build failure mentioning `undefined: Lookup`, `undefined: All`, `undefined: statuses`, `undefined: Status`.

- [ ] **Step 3: Write `status.go`**

Create `internal/status/status.go`:

```go
// Package status holds the table of HTTP status codes and their descriptions.
package status

import (
	"cmp"
	"slices"
)

// Status describes one HTTP status code.
type Status struct {
	Code        int
	Title       string
	Description string
	URL         string
}

// Lookup returns the status for code and whether it is known.
func Lookup(code int) (Status, bool) {
	s, ok := statuses[code]
	if !ok {
		return Status{}, false
	}
	s.Code = code
	return s, true
}

// All returns every known status sorted by code.
func All() []Status {
	all := make([]Status, 0, len(statuses))
	for code, s := range statuses {
		s.Code = code
		all = append(all, s)
	}
	slices.SortFunc(all, func(a, b Status) int { return cmp.Compare(a.Code, b.Code) })
	return all
}
```

- [ ] **Step 4: Generate `codes.go` mechanically from the old table**

The old table in `cmd/is.go` is regular: `"NNN": {` lines, then `code:`, `title:`, `description:`, `url:` fields. Convert it with this script, which drops the redundant `code:` field, switches keys to int, and exports the field names:

```bash
python3 - <<'PY'
import re, pathlib
src = pathlib.Path("cmd/is.go").read_text()
start = src.index("var statuses = map[string]Status{")
end = src.index("\n}\n", start) + len("\n}\n")
body = src[start:end]
body = body.replace("map[string]Status{", "map[int]Status{", 1)
body = re.sub(r'^\t"(\d+)": \{', r'\t\1: {', body, flags=re.M)
body = re.sub(r'^\t\tcode: .*\n', '', body, flags=re.M)
body = re.sub(r'^\t\ttitle: +', '\t\tTitle: ', body, flags=re.M)
body = re.sub(r'^\t\tdescription: +', '\t\tDescription: ', body, flags=re.M)
body = re.sub(r'^\t\turl: +', '\t\tURL: ', body, flags=re.M)
header = (
    "package status\n\n"
    "// statuses maps each known HTTP status code to its title, description, and\n"
    "// reference URL. Lookup and All fill in Code from the map key.\n"
)
pathlib.Path("internal/status/codes.go").write_text(header + body)
PY
gofmt -w internal/status/codes.go
grep -cE '^[[:space:]][0-9]+: \{' internal/status/codes.go
```

Expected: the final grep prints `55`. Inspect the first entry with `sed -n '1,12p' internal/status/codes.go`; it should read:

```go
package status

// statuses maps each known HTTP status code to its title, description, and
// reference URL. Lookup and All fill in Code from the map key.
var statuses = map[int]Status{
	100: {
		Title:       "Continue",
		Description: "The initial part of a request has been received and has not yet been rejected by the server. The server intends to send a final response after the request has been fully received and acted upon.",
		URL:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-100-continue",
	},
	101: {
```

- [ ] **Step 5: Run the tests to see the data bugs caught**

Run: `go test ./internal/status/`
Expected: FAIL. `TestTitlesMatchNetHTTP` reports `401: Title = "Not Authorized", want "Unauthorized" (net/http)` and `418: Title = "(Unused / I'm a teapot)", want "I'm a teapot" (net/http)`. `TestNewCodesPresent` reports all eight codes missing. The other three tests pass.

- [ ] **Step 6: Apply the data corrections**

```bash
python3 - <<'PY'
import pathlib
p = pathlib.Path("internal/status/codes.go")
s = p.read_text()
edits = [
    ('Title:       "Not Authorized",',
     'Title:       "Unauthorized",'),
    ('Title:       "(Unused / I\'m a teapot)",',
     'Title:       "I\'m a teapot",'),
    ('"The server refuses to brew coffee because it is, permanently, a teapot. This error is a reference to Hyper Text Coffee Pot Control Protocol defined in April Fools\' jokes in 1998 and 2014."',
     '"The server refuses to brew coffee because it is, permanently, a teapot. This is a reference to the Hyper Text Coffee Pot Control Protocol from the 1998 and 2014 April Fools\' RFCs. RFC 9110 lists 418 as unused."'),
    ('"The server has successfully fulfilled the request and that there is no additional content to send in the response content."',
     '"The server has successfully fulfilled the request and there is no additional content to send in the response content."'),
    ('"(305 was defined in a previous version of this specification and is now deprecated)"',
     '"305 was defined in a previous version of HTTP/1.1 and is now deprecated."'),
    ('"(306 was defined in a previous version of this specification, is no longer used, and the code is reserved.)"',
     '"306 was defined in a previous version of HTTP/1.1, is no longer used, and the code is reserved."'),
    ('"That the method received in the request-line is known by the origin server but not supported by the target resource."',
     '"The method received in the request-line is known by the origin server but not supported by the target resource."'),
    ('"Access to the target resource is no longer available at the origin server and that this condition is likely to be permanent."',
     '"Access to the target resource is no longer available at the origin server and this condition is likely to be permanent."'),
    ('"one or more conditions given in the request header fields evaluated to false when tested on the server."',
     '"One or more conditions given in the request header fields evaluated to false when tested on the server."'),
    ('"The policy for accessing the resource has not been met in the request."',
     '"The policy for accessing the resource has not been met in the request. This code is obsolete: RFC 2774 has been moved to Historic status."'),
    ('"https://html.spec.whatwg.org/multipage/semantics.html#early-hints"',
     '"https://www.rfc-editor.org/rfc/rfc8297"'),
    ('"https://httpwg.org/specs/rfc7725.html#n-451-unavailable-for-legal-reasons"',
     '"https://www.rfc-editor.org/rfc/rfc7725#section-3"'),
]
for old, new in edits:
    assert s.count(old) == 1, f"expected one match: {old[:50]}"
    s = s.replace(old, new)
p.write_text(s)
print("corrections applied")
PY
```

Expected: `corrections applied`.

- [ ] **Step 7: Add the eight missing codes**

Each new entry is inserted immediately before the existing entry with the next higher code so the file stays in ascending order.

```bash
python3 - <<'PY'
import pathlib
p = pathlib.Path("internal/status/codes.go")
s = p.read_text()

def entry(code, title, desc, url):
    return (f"\t{code}: {{\n"
            f"\t\tTitle:       \"{title}\",\n"
            f"\t\tDescription: \"{desc}\",\n"
            f"\t\tURL:         \"{url}\",\n"
            f"\t}},\n")

additions = [
    (103, entry(102, "Processing",
        "An interim response used to inform the client that the server has accepted the complete request but has not yet completed it. Deprecated by RFC 4918.",
        "https://www.rfc-editor.org/rfc/rfc2518#section-10.1")),
    (300, entry(207, "Multi-Status",
        "The message body that follows is by default an XML message and can contain a number of separate response codes, depending on how many sub-requests were made.",
        "https://www.rfc-editor.org/rfc/rfc4918#section-11.1")),
    (300, entry(208, "Already Reported",
        "The members of a DAV binding have already been enumerated in a preceding part of the multistatus response, and are not being included again.",
        "https://www.rfc-editor.org/rfc/rfc5842#section-7.1")),
    (300, entry(226, "IM Used",
        "The server has fulfilled a GET request for the resource, and the response is a representation of the result of one or more instance-manipulations applied to the current instance.",
        "https://www.rfc-editor.org/rfc/rfc3229#section-10.4.1")),
    (426, entry(423, "Locked",
        "The source or destination resource of a method is locked.",
        "https://www.rfc-editor.org/rfc/rfc4918#section-11.3")),
    (426, entry(424, "Failed Dependency",
        "The method could not be performed on the resource because the requested action depended on another action and that action failed.",
        "https://www.rfc-editor.org/rfc/rfc4918#section-11.4")),
    (426, entry(425, "Too Early",
        "The server is unwilling to risk processing a request that might be replayed.",
        "https://www.rfc-editor.org/rfc/rfc8470#section-5.2")),
    (429, entry(428, "Precondition Required",
        "The origin server requires the request to be conditional.",
        "https://www.rfc-editor.org/rfc/rfc6585#section-3")),
]
for before, text in additions:
    marker = f"\t{before}: {{\n"
    assert s.count(marker) == 1, f"marker for {before} not unique"
    s = s.replace(marker, text + marker)
p.write_text(s)
print("additions applied")
PY
gofmt -l internal/status/ ; grep -cE '^[[:space:]][0-9]+: \{' internal/status/codes.go
```

Expected: `additions applied`, `gofmt -l` prints nothing, grep prints `63`.

- [ ] **Step 8: Run the tests to verify they pass**

Run: `go test -race -cover -shuffle=on ./internal/status/`
Expected: `ok  	github.com/rpearce/httpwut/internal/status` with coverage 100.0%.

- [ ] **Step 9: Commit**

```bash
git add internal/status/
git commit -m "feat: add internal/status package with corrected and expanded data

Fixes the 401 title, several garbled descriptions, and the 103 and 451
links; adds 102, 207, 208, 226, 423, 424, 425, and 428.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 3: `internal/cli` package with testable commands

**Files:**
- Create: `internal/cli/root.go`
- Create: `internal/cli/is.go`
- Test: `internal/cli/cli_test.go`

**Interfaces:**
- Consumes: `status.Lookup(int) (status.Status, bool)`, `status.All() []status.Status` from Task 2; `browser.OpenURL(string) error` from `github.com/cli/browser`.
- Produces:
  - `func NewRootCmd(version string) *cobra.Command` (exported, used by Task 4).
  - `func newRootCmd(version string, openURL func(url string) error) *cobra.Command` (tests).
  - `func newIsCmd(openURL func(url string) error) *cobra.Command`.

- [ ] **Step 1: Write the failing tests**

Create `internal/cli/cli_test.go`:

```go
package cli

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"
)

// recorder is an openURL stub that records every URL and fails on failURL.
type recorder struct {
	urls    []string
	failURL string
	err     error
}

func (r *recorder) open(url string) error {
	r.urls = append(r.urls, url)
	if url == r.failURL {
		return r.err
	}
	return nil
}

// run executes the CLI with args and returns captured stdout, stderr, and
// the error from Execute. A nil args must become an empty slice, otherwise
// cobra falls back to os.Args.
func run(t *testing.T, rec *recorder, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	if args == nil {
		args = []string{}
	}
	var out, errOut bytes.Buffer
	root := newRootCmd("1.2.3", rec.open)
	root.SetOut(&out)
	root.SetErr(&errOut)
	root.SetArgs(args)
	err = root.Execute()
	return out.String(), errOut.String(), err
}

const verbose502 = "502 - Bad Gateway\n" +
	"The server, while acting as a gateway or proxy, received an invalid response from an inbound server it accessed while attempting to fulfill the request.\n" +
	"https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway\n"

func TestIs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     []string
		wantOut  string
		wantErr  string // substring of the error; empty means no error
		wantURLs []string
	}{
		{name: "code", args: []string{"is", "502"}, wantOut: "502 - Bad Gateway\n"},
		{name: "verbose", args: []string{"is", "502", "--verbose"}, wantOut: verbose502},
		{name: "short verbose", args: []string{"is", "502", "-v"}, wantOut: verbose502},
		{name: "leading zeros accepted", args: []string{"is", "0404"}, wantOut: "404 - Not Found\n"},
		{name: "non-numeric", args: []string{"is", "abc"}, wantErr: `invalid HTTP status code "abc"`},
		{name: "overflow", args: []string{"is", "99999999999999999999"}, wantErr: `invalid HTTP status code "99999999999999999999"`},
		{name: "unknown", args: []string{"is", "999"}, wantErr: "unknown HTTP status code 999"},
		{name: "negative", args: []string{"is", "--", "-1"}, wantErr: "unknown HTTP status code -1"},
		{name: "no args", args: []string{"is"}, wantErr: "accepts 1 arg(s), received 0"},
		{name: "two args", args: []string{"is", "200", "500"}, wantErr: "accepts 1 arg(s), received 2"},
		{name: "unknown flag", args: []string{"is", "404", "--nope"}, wantErr: "unknown flag: --nope"},
		{
			name:     "cats and dogs",
			args:     []string{"is", "404", "--cats", "--dogs"},
			wantOut:  "404 - Not Found\n",
			wantURLs: []string{"https://http.cat/404", "https://httpstatusdogs.com/404"},
		},
		{
			name:     "short flags",
			args:     []string{"is", "404", "-c", "-d"},
			wantOut:  "404 - Not Found\n",
			wantURLs: []string{"https://http.cat/404", "https://httpstatusdogs.com/404"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := &recorder{}
			out, errOut, err := run(t, rec, tt.args...)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
			}
			if out != tt.wantOut {
				t.Errorf("stdout = %q, want %q", out, tt.wantOut)
			}
			if errOut != "" {
				t.Errorf("stderr = %q, want empty", errOut)
			}
			if !slices.Equal(rec.urls, tt.wantURLs) {
				t.Errorf("opened URLs = %v, want %v", rec.urls, tt.wantURLs)
			}
		})
	}
}

func TestIsBrowserFailure(t *testing.T) {
	t.Parallel()

	boom := errors.New("boom")
	rec := &recorder{failURL: "https://http.cat/404", err: boom}
	out, _, err := run(t, rec, "is", "404", "--cats", "--dogs")
	if !errors.Is(err, boom) {
		t.Fatalf("error = %v, want wrapping %v", err, boom)
	}
	if !strings.Contains(err.Error(), "open http.cat: boom") {
		t.Errorf("error = %q, want containing %q", err.Error(), "open http.cat: boom")
	}
	if out != "404 - Not Found\n" {
		t.Errorf("stdout = %q, want the lookup line before the browser error", out)
	}
	want := []string{"https://http.cat/404", "https://httpstatusdogs.com/404"}
	if !slices.Equal(rec.urls, want) {
		t.Errorf("opened URLs = %v, want both attempted: %v", rec.urls, want)
	}
}

func TestRoot(t *testing.T) {
	t.Parallel()

	t.Run("version", func(t *testing.T) {
		t.Parallel()

		out, _, err := run(t, &recorder{}, "--version")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out != "httpwut version 1.2.3\n" {
			t.Errorf("stdout = %q, want %q", out, "httpwut version 1.2.3\n")
		}
	})

	t.Run("bare prints help", func(t *testing.T) {
		t.Parallel()

		out, _, err := run(t, &recorder{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "Usage:") || !strings.Contains(out, "  is ") {
			t.Errorf("stdout = %q, want help listing the is command", out)
		}
		if strings.Contains(out, "completion") {
			t.Errorf("stdout = %q, want the completion command hidden", out)
		}
	})

	t.Run("is help", func(t *testing.T) {
		t.Parallel()

		out, _, err := run(t, &recorder{}, "is", "-h")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "httpwut is <code> [flags]") {
			t.Errorf("stdout = %q, want the is usage line", out)
		}
	})

	t.Run("stray arg", func(t *testing.T) {
		t.Parallel()

		out, _, err := run(t, &recorder{}, "502")
		if err == nil || !strings.Contains(err.Error(), `unknown command "502" for "httpwut"`) {
			t.Fatalf("error = %v, want unknown command", err)
		}
		if out != "" {
			t.Errorf("stdout = %q, want empty", out)
		}
	})
}

func TestCompletion(t *testing.T) {
	t.Parallel()

	out, _, err := run(t, &recorder{}, "__complete", "is", "40")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) != 11 {
		t.Fatalf("got %d lines, want 10 codes plus a directive:\n%s", len(lines), out)
	}
	if lines[0] != "400\tBad Request" || lines[9] != "409\tConflict" {
		t.Errorf("first/last completions = %q, %q; want 400 and 409", lines[0], lines[9])
	}
	if lines[10] != ":4" {
		t.Errorf("directive line = %q, want %q (ShellCompDirectiveNoFileComp)", lines[10], ":4")
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./internal/cli/`
Expected: build failure `undefined: newRootCmd`.

- [ ] **Step 3: Write `root.go`**

Create `internal/cli/root.go`:

```go
// Package cli builds the httpwut command tree.
package cli

import (
	"github.com/cli/browser"
	"github.com/spf13/cobra"
)

// NewRootCmd builds the root command with the real browser opener.
func NewRootCmd(version string) *cobra.Command {
	return newRootCmd(version, browser.OpenURL)
}

func newRootCmd(version string, openURL func(url string) error) *cobra.Command {
	// No Args field on purpose: cobra's default rejects unknown subcommands
	// such as "httpwut 502", whereas cobra.NoArgs on a root without Run would
	// print help instead of an error.
	root := &cobra.Command{
		Use:           "httpwut",
		Short:         "Look up HTTP status code information",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}
	root.AddCommand(newIsCmd(openURL))
	return root
}
```

- [ ] **Step 4: Write `is.go`**

Create `internal/cli/is.go`:

```go
package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rpearce/httpwut/internal/status"
)

const (
	catsURL = "https://http.cat/"
	dogsURL = "https://httpstatusdogs.com/"
)

type isOptions struct {
	verbose bool
	cats    bool
	dogs    bool
}

func newIsCmd(openURL func(url string) error) *cobra.Command {
	var opts isOptions

	cmd := &cobra.Command{
		Use:               "is <code>",
		Short:             "Look up an HTTP status code",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeCodes,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIs(cmd, args[0], opts, openURL)
		},
	}

	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "Print status code description and URL")
	cmd.Flags().BoolVarP(&opts.cats, "cats", "c", false, "Open HTTP Status Cats webpage for status")
	cmd.Flags().BoolVarP(&opts.dogs, "dogs", "d", false, "Open HTTP Status Dogs webpage for status")

	return cmd
}

func runIs(cmd *cobra.Command, arg string, opts isOptions, openURL func(url string) error) error {
	code, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("invalid HTTP status code %q", arg)
	}

	s, ok := status.Lookup(code)
	if !ok {
		return fmt.Errorf("unknown HTTP status code %d", code)
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "%d - %s\n", s.Code, s.Title)
	if opts.verbose {
		fmt.Fprintln(out, s.Description)
		fmt.Fprintln(out, s.URL)
	}

	var errs []error
	if opts.cats {
		if openErr := openURL(catsURL + strconv.Itoa(s.Code)); openErr != nil {
			errs = append(errs, fmt.Errorf("open http.cat: %w", openErr))
		}
	}
	if opts.dogs {
		if openErr := openURL(dogsURL + strconv.Itoa(s.Code)); openErr != nil {
			errs = append(errs, fmt.Errorf("open httpstatusdogs.com: %w", openErr))
		}
	}
	return errors.Join(errs...)
}

// completeCodes offers every known status code whose digits start with what
// the user has typed so far.
func completeCodes(_ *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var comps []cobra.Completion
	for _, s := range status.All() {
		code := strconv.Itoa(s.Code)
		if strings.HasPrefix(code, toComplete) {
			comps = append(comps, cobra.CompletionWithDesc(code, s.Title))
		}
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}
```

- [ ] **Step 5: Run the tests to verify they pass**

Run: `go test -race -cover -shuffle=on ./internal/cli/ && go vet ./...`
Expected: `ok  	github.com/rpearce/httpwut/internal/cli` with coverage above 90%, no vet output. If `TestCompletion` fails on line count, print `out` and check which codes were returned; the expected set is exactly 400 through 409.

- [ ] **Step 6: Commit**

```bash
git add internal/cli/
git commit -m "feat: add internal/cli with testable, injectable commands

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 4: Wire `main.go`, add `--version`, delete `cmd/`

**Files:**
- Modify: `main.go` (full rewrite)
- Create: `main_test.go`
- Delete: `cmd/root.go`, `cmd/is.go`

**Interfaces:**
- Consumes: `cli.NewRootCmd(version string) *cobra.Command` from Task 3.
- Produces: `var version string` in package main (goreleaser ldflag target `main.version`); `func resolveVersion(ldflag string, info *debug.BuildInfo) string`.

- [ ] **Step 1: Write the failing test**

Create `main_test.go`:

```go
package main

import (
	"runtime/debug"
	"testing"
)

func TestResolveVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		ldflag string
		info   *debug.BuildInfo
		want   string
	}{
		{name: "ldflag wins", ldflag: "1.2.3", info: &debug.BuildInfo{Main: debug.Module{Version: "v9.9.9"}}, want: "1.2.3"},
		{name: "nil info", want: "dev"},
		{name: "empty module version", info: &debug.BuildInfo{}, want: "dev"},
		{name: "devel", info: &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, want: "dev"},
		{name: "stamped", info: &debug.BuildInfo{Main: debug.Module{Version: "v0.2.0"}}, want: "v0.2.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := resolveVersion(tt.ldflag, tt.info); got != tt.want {
				t.Errorf("resolveVersion(%q, %+v) = %q, want %q", tt.ldflag, tt.info, got, tt.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test .`
Expected: build failure `undefined: resolveVersion`.

- [ ] **Step 3: Rewrite `main.go`**

Replace the whole file with:

```go
// Command httpwut looks up HTTP status code information.
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/rpearce/httpwut/internal/cli"
)

// version is set by goreleaser via -ldflags "-X main.version=...".
var version string

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

// resolveVersion picks the ldflags value, then the module version stamped by
// the Go toolchain, then "dev".
func resolveVersion(ldflag string, info *debug.BuildInfo) string {
	if ldflag != "" {
		return ldflag
	}
	if info != nil && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
```

- [ ] **Step 4: Delete the old package and tidy**

Run: `git rm -rq cmd && go mod tidy && go build ./... && go vet ./...`
Expected: no output. `go.mod` is unchanged by tidy because cobra and cli/browser are still used.

- [ ] **Step 5: Run the whole test suite**

Run: `go test -race -cover -shuffle=on ./...`
Expected: `ok` for `github.com/rpearce/httpwut`, `.../internal/cli`, and `.../internal/status`.

- [ ] **Step 6: Smoke-test the binary end to end**

Run each line and check the result:

```bash
go run . is 502                      # 502 - Bad Gateway
go run . is 502 -v | wc -l           # 3
go run . is abc; echo "exit=$?"      # stderr: httpwut: invalid HTTP status code "abc"   exit=1
go run . is 999; echo "exit=$?"      # stderr: httpwut: unknown HTTP status code 999      exit=1
go run . is 200 500; echo "exit=$?"  # stderr: httpwut: accepts 1 arg(s), received 2     exit=1
go run . 502; echo "exit=$?"         # stderr: httpwut: unknown command "502" for "httpwut"  exit=1
go run . --version                   # httpwut version dev  OR  httpwut version v0.1.1-0.<timestamp>-<hash>
go run . -h | grep -c completion     # 0
```

Expected: as in the trailing comments. No error case prints a `Usage:` block. The version line is `dev` when the toolchain did not stamp VCS info and a `v0.1.1-0.…` pseudo-version when it did; both are correct.

- [ ] **Step 7: Commit**

```bash
git add main.go main_test.go go.mod go.sum
git commit -m "feat: wire main to internal/cli and add --version

Replaces the cobra-cli cmd/ package. Errors now print as
\"httpwut: <message>\" without the usage block.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 5: golangci-lint configuration

**Files:**
- Create: `.golangci.yml`

**Interfaces:**
- Consumes: all Go code from Tasks 2 through 4.
- Produces: a lint config that CI (Task 7) and the `mise run lint` task use.

- [ ] **Step 1: Create `.golangci.yml`**

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
      extra-rules: true
    goimports:
      local-prefixes:
        - github.com/rpearce/httpwut
```

- [ ] **Step 2: Run the formatter check**

Run: `golangci-lint fmt --diff`
Expected: no output. If a diff is printed, run `golangci-lint fmt` to apply it and re-run the check.

- [ ] **Step 3: Run the linter**

Run: `golangci-lint run`
Expected: `0 issues.`

Known possible finding and its fix: if revive's `package-comments` rule rejects the `// Command httpwut ...` comment in `main.go` with "package comment should be of the form", change that line to `// Package main implements httpwut, a CLI that looks up HTTP status code information.` and re-run. Any other finding is a defect in the code from Tasks 2 through 4; fix it at the source, keep the tests green, and note the change in the commit body.

- [ ] **Step 4: Run the tests once more**

Run: `go test -race -cover -shuffle=on ./...`
Expected: all `ok`.

- [ ] **Step 5: Commit**

```bash
git add .golangci.yml main.go internal/
git commit -m "chore: add golangci-lint v2 config

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 6: goreleaser v2 configuration

**Files:**
- Modify: `.goreleaser.yaml` (full rewrite)

**Interfaces:**
- Consumes: `main.version` from Task 4.
- Produces: a config that `goreleaser check` accepts and that the release workflow (Task 7) runs.

- [ ] **Step 1: Rewrite `.goreleaser.yaml`**

Replace the whole file with:

```yaml
version: 2

before:
  hooks:
    - go mod tidy

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
      - -s -w -X main.version={{ .Version }}

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

- [ ] **Step 2: Validate the config**

Run: `goreleaser check`
Expected: `1 configuration file(s) validated` with no deprecation warnings.

- [ ] **Step 3: Build a snapshot**

Run: `goreleaser release --snapshot --clean 2>&1 | tail -5 && ls dist/*.tar.gz dist/*.zip dist/checksums.txt`
Expected: the build succeeds and `dist/` contains exactly these archives plus `checksums.txt`:

```
dist/httpwut_0.1.2-next_darwin_amd64.tar.gz
dist/httpwut_0.1.2-next_darwin_arm64.tar.gz
dist/httpwut_0.1.2-next_linux_amd64.tar.gz
dist/httpwut_0.1.2-next_linux_arm64.tar.gz
dist/httpwut_0.1.2-next_windows_amd64.zip
dist/httpwut_0.1.2-next_windows_arm64.zip
```

- [ ] **Step 4: Verify the ldflag reached the binary**

Run: `"$(find dist -type f -name httpwut -path '*darwin_arm64*')" --version`
Expected: `httpwut version 0.1.2-next`.

- [ ] **Step 5: Confirm nothing untracked leaked**

Run: `git status --short`
Expected: only ` M .goreleaser.yaml`. `dist/` is gitignored.

- [ ] **Step 6: Commit**

```bash
git add .goreleaser.yaml
git commit -m "build: rewrite goreleaser config for v2

Drops 386 builds, adds -trimpath, a reproducible mod_timestamp, and a
version ldflag. Archive names now follow the goreleaser default.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 7: GitHub Actions and Dependabot

**Files:**
- Create: `.github/workflows/ci.yml`
- Create: `.github/workflows/release.yml`
- Create: `.github/dependabot.yml`

**Interfaces:**
- Consumes: `.golangci.yml` (Task 5), `.goreleaser.yaml` (Task 6), `go.mod` go directive (Task 1).
- Produces: CI on push and pull request; releases on `v*` tags.

- [ ] **Step 1: Create `.github/workflows/ci.yml`**

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

permissions:
  contents: read

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
          check-latest: true
      - run: go mod tidy -diff
      - run: go vet ./...
      - run: go test -race -cover -shuffle=on ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
          check-latest: true
      - uses: golangci/golangci-lint-action@v9
        with:
          version: v2.13.2

  vuln:
    runs-on: ubuntu-latest
    steps:
      - uses: golang/govulncheck-action@v1
        with:
          go-version-file: go.mod
```

- [ ] **Step 2: Create `.github/workflows/release.yml`**

```yaml
name: Release

on:
  push:
    tags: ["v*"]

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
      - uses: goreleaser/goreleaser-action@v7
        with:
          distribution: goreleaser
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

- [ ] **Step 3: Create `.github/dependabot.yml`**

```yaml
version: 2

updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly

  - package-ecosystem: github-actions
    directory: /
    schedule:
      interval: weekly
```

- [ ] **Step 4: Validate the workflow files**

Run: `mise exec actionlint@latest -- actionlint`
Expected: no output, exit 0. actionlint checks syntax, job and step shapes, and expression references. It does not check that the action versions exist; those were verified against the GitHub releases API on 2026-09-22 when the spec was written.

If mise cannot install actionlint, fall back to a YAML parse: `for f in .github/workflows/*.yml .github/dependabot.yml; do ruby -ryaml -e 'YAML.load_file(ARGV[0]); puts "ok #{ARGV[0]}"' "$f"; done`.

- [ ] **Step 5: Run the same commands CI will run, locally**

Run: `go mod tidy -diff && go vet ./... && go test -race -cover -shuffle=on ./... && golangci-lint run && govulncheck ./...`
Expected: tidy prints nothing, vet prints nothing, tests `ok`, `0 issues.`, and govulncheck ends with `No vulnerabilities found.`

- [ ] **Step 6: Commit**

```bash
git add .github/
git commit -m "ci: add test, lint, vuln, and release workflows with dependabot

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 8: README refresh

**Files:**
- Modify: `README.md` (full rewrite)

**Interfaces:**
- Consumes: actual program output from Task 4 and mise tasks from Task 1.
- Produces: documentation that matches the binary.

- [ ] **Step 1: Capture the real help text**

Run: `go run . is -h`
Expected:

```
Look up an HTTP status code

Usage:
  httpwut is <code> [flags]

Flags:
  -c, --cats      Open HTTP Status Cats webpage for status
  -d, --dogs      Open HTTP Status Dogs webpage for status
  -h, --help      help for is
  -v, --verbose   Print status code description and URL
```

If the output differs in any character, use the real output in the next step.

- [ ] **Step 2: Rewrite `README.md`**

Replace the whole file with (the `is -h` block must match Step 1 exactly):

````markdown
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
````

- [ ] **Step 3: Verify the examples against the binary**

Run:

```bash
diff <(go run . is -h; echo) <(sed -n '/^λ httpwut is -h$/,/^λ httpwut is 502$/p' README.md | sed '1d;$d') && echo HELP_MATCHES
diff <(go run . is 502 --verbose) <(sed -n '/^λ httpwut is 502 --verbose$/,/^```$/p' README.md | sed '1d;$d') && echo VERBOSE_MATCHES
go run . completion --help | head -1
mise tasks ls
```

Expected: `HELP_MATCHES`, `VERBOSE_MATCHES`, a first help line for the completion command, and a task list showing `fmt`, `lint`, `snapshot`, `test`, `vuln`.

- [ ] **Step 4: Commit**

```bash
git add README.md
git commit -m "docs: refresh README for the modernized CLI

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 9: Definition-of-done run

**Files:** none modified. This task only verifies.

**Interfaces:**
- Consumes: everything above.
- Produces: evidence that the spec's section 13 holds.

- [ ] **Step 1: Run the full verification list from the spec**

```bash
go build ./... && echo BUILD_OK
go vet ./... && echo VET_OK
go test -race -cover -shuffle=on ./...
golangci-lint run
golangci-lint fmt --diff && echo FMT_OK
govulncheck ./...
goreleaser check
goreleaser release --snapshot --clean >/dev/null 2>&1 && ls dist/*.tar.gz dist/*.zip | wc -l
go run . --version
go run . is 502
go run . is 502 -v
go run . is abc; echo "exit=$?"
go run . is 999; echo "exit=$?"
go run . is 200 500; echo "exit=$?"
go run . 502; echo "exit=$?"
git status --short
git log --oneline main..HEAD
```

Expected, in order: `BUILD_OK`; `VET_OK`; three `ok` lines; `0 issues.`; `FMT_OK`; `No vulnerabilities found.`; `1 configuration file(s) validated`; `6`; a version line; `502 - Bad Gateway`; three lines; `httpwut: invalid HTTP status code "abc"` then `exit=1`; `httpwut: unknown HTTP status code 999` then `exit=1`; `httpwut: accepts 1 arg(s), received 2` then `exit=1`; `httpwut: unknown command "502" for "httpwut"` then `exit=1`; an empty status; and eleven commits: two for the spec, one for the plan, and one per Task 1 through 8. Extra commits from fixes made during Task 5 are fine.

No error case may print a `Usage:` block.

- [ ] **Step 2: Report**

Summarize the verification output verbatim for the reviewer. Do not push, tag, or open a pull request; those are the owner's calls.
