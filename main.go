// Command httpwut looks up HTTP status code information.
package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/rpearce/httpwut/internal/cli"
)

// version is set by goreleaser via -ldflags "-X main.version=...".
var version string

func main() {
	if err := cli.NewRootCmd(buildVersion()).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, errorLine(err))
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

// errorLine formats err for stderr. Each error in an errors.Join result gets
// its own prefixed line. Cobra's unknown-command suggestions end in a newline,
// which Fprintln would double, so trailing newlines are trimmed.
//
// errors.As finds the first joined error in the chain, so wrapping a Join with
// fmt.Errorf("ctx: %w", ...) would drop "ctx: ". runIs returns its Join
// unwrapped.
func errorLine(err error) string {
	var joined interface{ Unwrap() []error }
	if errors.As(err, &joined) {
		errs := joined.Unwrap()
		lines := make([]string, 0, len(errs))
		for _, e := range errs {
			lines = append(lines, errorLine(e))
		}
		return strings.Join(lines, "\n")
	}
	return "httpwut: " + strings.TrimRight(err.Error(), "\n")
}
