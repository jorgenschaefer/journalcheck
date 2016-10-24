package filter

import "github.com/coreos/go-systemd/sdjournal"

type Filter interface {
	Matches(entry *sdjournal.JournalEntry) bool
}
