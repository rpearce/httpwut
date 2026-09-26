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
