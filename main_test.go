package main

import (
	"errors"
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

func TestErrorLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		msg  string
		want string
	}{
		{"plain", "boom", "httpwut: boom"},
		{
			"trailing newline is trimmed, inner newlines kept",
			"unknown command \"si\" for \"httpwut\"\n\nDid you mean this?\n\tis\n",
			"httpwut: unknown command \"si\" for \"httpwut\"\n\nDid you mean this?\n\tis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := errorLine(errors.New(tt.msg)); got != tt.want {
				t.Errorf("errorLine(%q) = %q, want %q", tt.msg, got, tt.want)
			}
		})
	}
}

func TestErrorLineJoined(t *testing.T) {
	t.Parallel()

	// The second message ends in a newline to check trimming per error.
	msgs := []string{"open http.cat: boom", "open httpstatusdogs.com: boom\n"}
	err := errors.Join(errors.New(msgs[0]), errors.New(msgs[1]))
	want := "httpwut: open http.cat: boom\nhttpwut: open httpstatusdogs.com: boom"
	if got := errorLine(err); got != want {
		t.Errorf("errorLine(joined) = %q, want %q", got, want)
	}
}
