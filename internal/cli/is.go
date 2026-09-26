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
