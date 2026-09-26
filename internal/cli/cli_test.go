package cli

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"
)

// recorder is an openURL stub that records every URL and fails on failURL,
// or on every URL when failAll is set.
type recorder struct {
	urls    []string
	failURL string
	failAll bool
	err     error
}

func (r *recorder) open(url string) error {
	r.urls = append(r.urls, url)
	if r.failAll || url == r.failURL {
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
		{
			name:     "cats only",
			args:     []string{"is", "404", "--cats"},
			wantOut:  "404 - Not Found\n",
			wantURLs: []string{"https://http.cat/404"},
		},
		{
			name:     "dogs only",
			args:     []string{"is", "404", "--dogs"},
			wantOut:  "404 - Not Found\n",
			wantURLs: []string{"https://httpstatusdogs.com/404"},
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

	tests := []struct {
		name    string
		rec     *recorder
		wantErr []string // substrings that must all appear in the error
	}{
		{
			name:    "cats fails",
			rec:     &recorder{failURL: "https://http.cat/404", err: boom},
			wantErr: []string{"open http.cat: boom"},
		},
		{
			name:    "dogs fails",
			rec:     &recorder{failURL: "https://httpstatusdogs.com/404", err: boom},
			wantErr: []string{"open httpstatusdogs.com: boom"},
		},
		{
			name:    "both fail",
			rec:     &recorder{failAll: true, err: boom},
			wantErr: []string{"open http.cat: boom", "open httpstatusdogs.com: boom"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			out, _, err := run(t, tt.rec, "is", "404", "--cats", "--dogs")
			if !errors.Is(err, boom) {
				t.Fatalf("error = %v, want wrapping %v", err, boom)
			}
			for _, want := range tt.wantErr {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error = %q, want containing %q", err.Error(), want)
				}
			}
			if out != "404 - Not Found\n" {
				t.Errorf("stdout = %q, want the lookup line before the browser error", out)
			}
			want := []string{"https://http.cat/404", "https://httpstatusdogs.com/404"}
			if !slices.Equal(tt.rec.urls, want) {
				t.Errorf("opened URLs = %v, want both attempted: %v", tt.rec.urls, want)
			}
		})
	}
}

func TestRoot(t *testing.T) {
	t.Parallel()

	t.Run("version", func(t *testing.T) {
		t.Parallel()

		out, errOut, err := run(t, &recorder{}, "--version")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out != "httpwut version 1.2.3\n" {
			t.Errorf("stdout = %q, want %q", out, "httpwut version 1.2.3\n")
		}
		if errOut != "" {
			t.Errorf("stderr = %q, want empty", errOut)
		}
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()

		// Root's -v is cobra's version shorthand; is -v remains --verbose.
		out, errOut, err := run(t, &recorder{}, "-v")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out != "httpwut version 1.2.3\n" {
			t.Errorf("stdout = %q, want %q", out, "httpwut version 1.2.3\n")
		}
		if errOut != "" {
			t.Errorf("stderr = %q, want empty", errOut)
		}
	})

	t.Run("bare prints help", func(t *testing.T) {
		t.Parallel()

		out, errOut, err := run(t, &recorder{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "Usage:") || !strings.Contains(out, "  is ") {
			t.Errorf("stdout = %q, want help listing the is command", out)
		}
		if strings.Contains(out, "completion") {
			t.Errorf("stdout = %q, want the completion command hidden", out)
		}
		if errOut != "" {
			t.Errorf("stderr = %q, want empty", errOut)
		}
	})

	t.Run("is help", func(t *testing.T) {
		t.Parallel()

		out, errOut, err := run(t, &recorder{}, "is", "-h")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "httpwut is <code> [flags]") {
			t.Errorf("stdout = %q, want the is usage line", out)
		}
		if errOut != "" {
			t.Errorf("stderr = %q, want empty", errOut)
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
