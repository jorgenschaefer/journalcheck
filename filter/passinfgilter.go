package filter

import "github.com/coreos/go-systemd/sdjournal"

type PassingFilter struct {
}

func NewPassingFilter() Filter {
	return new(PassingFilter)
}

func (f *PassingFilter) Matches(entry *sdjournal.JournalEntry) bool {
	return false
}
