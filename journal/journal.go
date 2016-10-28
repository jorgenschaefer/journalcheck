package journal

import (
	"fmt"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
)

type Entry sdjournal.JournalEntry
type Journal sdjournal.Journal

func New() (*Journal, error) {
	j, err := sdjournal.NewJournal()
	return (*Journal)(j), err
}

func (j *Journal) SeekCursor(cursor string) error {
	sdj := (*sdjournal.Journal)(j)
	if err := sdj.SeekCursor(cursor); err != nil {
		return err
	}
	// Move to the position of the cursor, else we'd see the last
	// entry again
	_, err := sdj.Next()
	return err
}

func (j *Journal) SeekLast(last int) error {
	sdj := (*sdjournal.Journal)(j)
	if err := sdj.SeekTail(); err != nil {
		return err
	}
	_, err := sdj.PreviousSkip(uint64(last) + 1)
	return err
}

func (j *Journal) Next() (entry *Entry, err error) {
	sdj := (*sdjournal.Journal)(j)
	n, err := sdj.Next()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	e, err := sdj.GetEntry()
	if err != nil {
		return nil, err
	}
	return (*Entry)(e), nil
}

func (e *Entry) MatchString() string {
	identifier := getFirst(e.Fields, "_SYSTEMD_UNIT", "_COMM", "SYSLOG_IDENTIFIER")
	message := e.Fields["MESSAGE"]
	return fmt.Sprintf("%s: %s", identifier, message)
}

func (e *Entry) ShortString() string {
	logtime := time.Unix(int64(e.RealtimeTimestamp/1000/1000), 0).Format(time.Stamp)
	host := getFirst(e.Fields, "_HOSTNAME")
	identifier := getFirst(e.Fields, "_SYSTEMD_UNIT", "_COMM", "SYSLOG_IDENTIFIER")
	pid := getFirst(e.Fields, "SYSLOG_PID", "_PID")
	message := e.Fields["MESSAGE"]
	return fmt.Sprintf("%s %s %s[%s]: %s", logtime, host, identifier, pid, message)
}

func (e *Entry) VerboseString() string {
	panic("VerboseString of journal entries not implemented yet")
}

func getFirst(entries map[string]string, fields ...string) string {
	for _, field := range fields {
		value, ok := entries[field]
		if ok {
			return value
		}
	}
	return "<unset>"
}
