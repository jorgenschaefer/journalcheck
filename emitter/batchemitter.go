package emitter

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/coreos/go-systemd/sdjournal"

	"github.com/jorgenschaefer/journalcheck/buffer"
	"github.com/jorgenschaefer/journalcheck/filter"
)

type BatchEmitter struct {
	cursorFile  *string
	filter      filter.Filter
	maxLen      int
	maxDuration time.Duration
	timeout     time.Duration
	sender      func([]*sdjournal.JournalEntry)
}

func NewBatchEmitter() *BatchEmitter {
	return &BatchEmitter{
		maxLen:      100,
		maxDuration: time.Hour,
		timeout:     time.Duration(15) * time.Minute,
	}
}

func (e *BatchEmitter) SetCursorFile(cursorFile string) {
	e.cursorFile = &cursorFile
}

func (e *BatchEmitter) SetFilter(filter filter.Filter) {
	e.filter = filter
}

func (e *BatchEmitter) Consume(entries chan *sdjournal.JournalEntry) {
	var buffer = buffer.NewBuffer(e.maxLen)
	var firstEntry time.Time
	var lastCursor string = ""
	var haveCursor bool = false

	for {
		entry, timeout := getEntry(entries, e.timeout)
		if !timeout {
			if e.filter == nil || !e.filter.Matches(entry) {
				if buffer.IsEmpty() {
					firstEntry = time.Now()
				}
				buffer, err := buffer.Append(entry)
				if err != nil {
					log.Fatal(err)
				}
			}
			lastCursor = entry.Cursor
			haveCursor = true
		}
		if timeout || buffer.IsFull() {
			e.sender(buffer)
			buffer.Clear()
		}
		if len(buffer) == 0 && haveCursor {
			e.storeCursor(lastCursor)
			haveCursor = false
		}
	}
}

func (e *BatchEmitter) storeCursor(cursor string) {
	if e.cursorFile != nil {
		if err := ioutil.WriteFile(*e.cursorFile, []byte(cursor), 0600); err != nil {
			log.Fatal(err)
		}
	}
}

func getEntry(entries chan *sdjournal.JournalEntry, timeout time.Duration) (*sdjournal.JournalEntry, bool) {
	select {
	case entry := <-entries:
		return entry, false
	case <-time.After(timeout):
		return nil, true
	}
}
