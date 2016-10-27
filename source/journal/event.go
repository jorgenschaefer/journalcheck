package journal

import (
	"fmt"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/jorgenschaefer/journalcheck/event"
)

type JournalEvent struct {
	entry *sdjournal.JournalEntry
}

func NewEvent(entry *sdjournal.JournalEntry) event.Event {
	return &JournalEvent{
		entry: entry,
	}
}

func (e *JournalEvent) MatchString() string {
	identifier := getField(e.entry, "_SYSTEMD_UNIT", "_COMM", "SYSLOG_IDENTIFIER")
	message := e.entry.Fields["MESSAGE"]
	return fmt.Sprintf("%s: %s", identifier, message)
}

func (e *JournalEvent) ShortString() string {
	logtime := time.Unix(int64(e.entry.RealtimeTimestamp/1000/1000), 0).Format(time.Stamp)
	host := getField(e.entry, "_HOSTNAME")
	identifier := getField(e.entry, "_SYSTEMD_UNIT", "_COMM", "SYSLOG_IDENTIFIER")
	pid := getField(e.entry, "SYSLOG_PID", "_PID")
	message := e.entry.Fields["MESSAGE"]
	return fmt.Sprintf("%s %s %s[%s]: %s", logtime, host, identifier, pid, message)
}

func (e *JournalEvent) LongString() string {
	panic("LongString of journal entries not implemented yet")
}

func (e *JournalEvent) Sent() {
}

func getField(entry *sdjournal.JournalEntry, fields ...string) string {
	for _, field := range fields {
		value, ok := entry.Fields[field]
		if ok {
			return value
		}
	}
	return "<unset>"
}
