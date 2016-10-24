package formatter

import (
	"fmt"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
)

func Match(entry *sdjournal.JournalEntry) string {
	identifier := getField(entry, "_SYSTEMD_UNIT", "_COMM", "SYSLOG_IDENTIFIER")
	message := entry.Fields["MESSAGE"]
	return fmt.Sprintf("%s: %s", identifier, message)
}

func Syslog(entry *sdjournal.JournalEntry) string {
	logtime := time.Unix(int64(entry.RealtimeTimestamp/1000/1000), 0).Format(time.Stamp)
	identifier := getField(entry, "_SYSTEMD_UNIT", "_COMM", "SYSLOG_IDENTIFIER")
	pid := getField(entry, "SYSLOG_PID", "_PID")
	message := entry.Fields["MESSAGE"]
	return fmt.Sprintf("%s %s[%s]: %s", logtime, identifier, pid, message)
}

func getField(entry *sdjournal.JournalEntry, fields ...string) string {
	for _, field := range fields {
		value, ok := entry.Fields[field]
		if ok {
			return value
		}
	}
	return ""
}
