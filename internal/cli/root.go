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
