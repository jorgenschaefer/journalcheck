// Package journal provides a simplified journal API. The sdjournal
// library is quite lowlevel, the API provided here is a bit more
// high-level and more intuitive to use.
package journal

import (
	"fmt"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
)

// Journal is a handle for a journal to talk to.
type Journal sdjournal.Journal

// Entry is a single entry from the journal.
type Entry sdjournal.JournalEntry

// New returns a handle to the journal the current user can read. For
// root, this is the system journal. For other users, this is usually
// just their own journal.
func New() (*Journal, error) {
	j, err := sdjournal.NewJournal()
	return (*Journal)(j), err
}

// SeekCursor moves the current cursor past the element pointed to by
// the cursor name, so further reading proceeds after this last
// element.
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

// SeekLast moves the cursor so that further processing will result in
// `last` entries until the end of the journal (provided no new
// entries are added in the meantime).
func (j *Journal) SeekLast(last int) error {
	sdj := (*sdjournal.Journal)(j)
	if err := sdj.SeekTail(); err != nil {
		return err
	}
	_, err := sdj.PreviousSkip(uint64(last) + 1)
	return err
}

// Next returns the next entry. This entry is nil if the cursor is at
// the end of the journal.
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

// MatchString returns the entry formatted for matching. This includes
// the systemd unit name and the message only. If the systemd unit
// name is not available, this falls back to the command, and if that
// is also not available, the syslog identifier.
func (e *Entry) MatchString() string {
	identifier := getFirst(e.Fields, "SYSLOG_IDENTIFIER", "_SYSTEMD_UNIT", "_COMM")
	message := e.Fields["MESSAGE"]
	return fmt.Sprintf("%s: %s", identifier, message)
}

// ShortString returns the entry as a single line, very similar to a
// typical line from syslog.
func (e *Entry) ShortString() string {
	logtime := time.Unix(int64(e.RealtimeTimestamp/1000/1000), 0).Format(time.Stamp)
	host := getFirst(e.Fields, "_HOSTNAME")
	identifier := getFirst(e.Fields, "SYSLOG_IDENTIFIER", "_SYSTEMD_UNIT", "_COMM")
	pid := getFirst(e.Fields, "SYSLOG_PID", "_PID")
	message := e.Fields["MESSAGE"]
	return fmt.Sprintf("%s %s %s[%s]: %s", logtime, host, identifier, pid, message)
}

// VerboseString returns the entry formatted like journalctl's verbose
// output format, with one key-value-pair per line.
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
