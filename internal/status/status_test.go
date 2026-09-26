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
