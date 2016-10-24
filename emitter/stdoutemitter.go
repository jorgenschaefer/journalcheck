package emitter

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/coreos/go-systemd/sdjournal"

	"github.com/jorgenschaefer/journalcheck/filter"
	"github.com/jorgenschaefer/journalcheck/formatter"
)

type StdoutEmitter struct {
	cursorFile *string
	filter     filter.Filter
}

func NewStdoutEmitter() *StdoutEmitter {
	return new(StdoutEmitter)
}

func (e *StdoutEmitter) SetCursorFile(cursorFile string) {
	e.cursorFile = &cursorFile
}

func (e *StdoutEmitter) SetFilter(filter filter.Filter) {
	e.filter = filter
}

func (e *StdoutEmitter) Consume(entries chan *sdjournal.JournalEntry) {
	for entry := range entries {
		if e.shouldShow(entry) {
			fmt.Println(formatter.Syslog(entry))
		}
		e.storeCursor(entry)
	}
}

func (e *StdoutEmitter) shouldShow(entry *sdjournal.JournalEntry) bool {
	return e.filter == nil || !e.filter.Matches(entry)
}

func (e *StdoutEmitter) storeCursor(entry *sdjournal.JournalEntry) {
	if e.cursorFile != nil {
		if err := ioutil.WriteFile(*e.cursorFile, []byte(entry.Cursor), 0600); err != nil {
			log.Fatal(err)
		}
	}
}
