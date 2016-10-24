package emitter

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/coreos/go-systemd/sdjournal"

	"github.com/jorgenschaefer/journalcheck/filter"
)

type EmailEmitter struct {
	cursorFile     *string
	filter         filter.Filter
	recipient      string
	MaxLen         int
	MaxDuration    time.Duration
	TimeoutMinutes int
}

func NewEmailEmitter(recipient string) *EmailEmitter {
	emitter := new(EmailEmitter)
	emitter.recipient = recipient
	emitter.MaxLen = 1000
	emitter.MaxDuration = time.Hour
	emitter.TimeoutMinutes = 60
	return emitter
}

func (e *EmailEmitter) SetCursorFile(cursorFile string) {
	e.cursorFile = &cursorFile
}

func (e *EmailEmitter) SetFilter(filter filter.Filter) {
	e.filter = filter
}

func (e *EmailEmitter) Consume(entries chan *sdjournal.JournalEntry) {
	var sender *Sender = MakeSender(e.recipient, e.MaxLen, e.MaxDuration)
	var lastCursor *string = nil

	for {
		entry, timeout := getEntry(entries, e.TimeoutMinutes)
		if !timeout {
			if e.filter == nil || !e.filter.Matches(entry) {
				sender.Add(entry)
			}
			lastCursor = &entry.Cursor
		}
		if timeout || sender.ShouldSend() {
			sender.Send()
		}
		if sender.IsEmpty() && lastCursor != nil {
			e.storeCursor(*lastCursor)
			lastCursor = nil
		}
	}
}

func (e *EmailEmitter) storeCursor(cursor string) {
	if e.cursorFile != nil {
		if err := ioutil.WriteFile(*e.cursorFile, []byte(cursor), 0600); err != nil {
			log.Fatal(err)
		}
	}
}

func getEntry(entries chan *sdjournal.JournalEntry, timeoutMinutes int) (*sdjournal.JournalEntry, bool) {
	select {
	case entry := <-entries:
		return entry, false
	case <-time.After(time.Minute * time.Duration(timeoutMinutes)):
		return nil, true
	}
}

type Sender struct {
	entries     []*sdjournal.JournalEntry
	first       time.Time
	maxDuration time.Duration
	recipient   string
}

func MakeSender(recipient string, maxLen int, maxDuration time.Duration) *Sender {
	sender := new(Sender)
	sender.entries = make([]*sdjournal.JournalEntry, 0, maxLen)
	sender.maxDuration = maxDuration
	sender.recipient = recipient
	return sender
}

func (s *Sender) Add(entry *sdjournal.JournalEntry) {
	if s.IsEmpty() {
		s.first = time.Now()
	}
	s.entries = append(s.entries, entry)
}

func (s *Sender) ShouldSend() bool {
	if len(s.entries) == 0 {
		return false
	}
	if len(s.entries) == cap(s.entries) {
		return true
	}
	if time.Now().Sub(s.first) > s.maxDuration {
		return true
	}
	return false
}

func (s *Sender) Send() {
	if len(s.entries) == 0 {
		return
	}
	log.Printf("Sending %d entries\n", len(s.entries))
	// sendmail := exec.Command("/usr/sbin/sendmail", s.recipient)
	// stdin, err := sendmail.StdinPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := sendmail.Start(); err != nil {
	// 	log.Fatal(err)
	// }

	// mail := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n",
	// 	"journalcheck", s.recipient, "Journal output")
	// stdin.Write([]byte(mail))
	// for _, entry := range s.entries {
	// 	stdin.Write([]byte(formatter.Syslog(entry)))
	// 	stdin.Write([]byte("\n"))
	// }
	// stdin.Close()
	// sendmail.Wait()

	s.entries = make([]*sdjournal.JournalEntry, 0, cap(s.entries))
}

func (s *Sender) IsEmpty() bool {
	return len(s.entries) == 0
}
