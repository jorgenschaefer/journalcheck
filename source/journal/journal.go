package journal

import (
	"io/ioutil"
	"log"

	"github.com/coreos/go-systemd/sdjournal"

	"github.com/jorgenschaefer/journalcheck/event"
	"github.com/jorgenschaefer/journalcheck/source"
)

type journalSource struct {
	journal   *sdjournal.Journal
	terminate bool
	filter    Filter
}

func NewDelimitedSource(entryCount int, filterfile *string) source.Source {
	source := filteredSource(filterfile)
	source.SeekLast(entryCount)
	return source
}

func NewCursorSource(cursorfile string, filterfile *string) source.Source {
	source := filteredSource(filterfile)
	cursor, err := ioutil.ReadFile(cursorfile)
	if err != nil {
		log.Fatal(err)
	}
	source.SeekCursor(string(cursor))
	return source
}

func filteredSource(filterfile *string) *journalSource {
	journal, err := sdjournal.NewJournal()
	if err != nil {
		log.Fatal(err)
	}
	var filter Filter
	if filterfile != nil {
		filter, err = NewFilter(*filterfile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		filter = NewEmptyFilter()
	}
	return &journalSource{
		journal:   journal,
		terminate: true,
		filter:    filter,
	}
}

func (s *journalSource) SeekCursor(cursor string) {
	if err := s.journal.SeekCursor(cursor); err != nil {
		log.Fatal(err)
	}
	// Move to the position of the cursor, else we'd see the last
	// entry again
	_, err := s.journal.Next()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *journalSource) SeekLast(rewind int) {
	if err := s.journal.SeekTail(); err != nil {
		log.Fatal(err)
	}
	_, err := s.journal.PreviousSkip(uint64(rewind) + 1)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *journalSource) Emit(eventstream chan event.Event) {
	for {
		n, err := s.journal.Next()
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			if s.terminate {
				close(eventstream)
				return
			} else {
				s.journal.Wait(sdjournal.IndefiniteWait)
				continue
			}
		}
		entry, err := s.journal.GetEntry()
		if err != nil {
			log.Fatal(err)
		}
		event := NewEvent(entry)
		if !s.filter.Matches(event) {
			eventstream <- event
		}
	}
}
