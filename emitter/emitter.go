package emitter

import (
	"github.com/coreos/go-systemd/sdjournal"
	"github.com/jorgenschaefer/journalcheck/filter"
)

type Emitter interface {
	SetCursorFile(string)
	SetFilter(filter.Filter)
	Consume(chan *sdjournal.JournalEntry)
}
