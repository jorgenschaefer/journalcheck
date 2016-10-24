package formatter

import "testing"
import "github.com/coreos/go-systemd/sdjournal"

func TestFormatLine(t *testing.T) {
	entry := new(sdjournal.JournalEntry)
	entry.Fields = map[string]string{
		"MESSAGE":           "Hello, World",
		"SYSLOG_IDENTIFIER": "program",
		"SYSLOG_PID":        "23",
		"_PID":              "42",
		"_SYSTEMD_UNIT":     "some.service",
		"_TRANSPORT":        "driver",
	}
	expectedSyslog := "program[23]: Hello, World"
	expectedJournal := "some.service: Hello, World"

	testCases := map[string]string{
		"audit":   expectedJournal,
		"driver":  expectedJournal,
		"syslog":  expectedSyslog,
		"journal": expectedJournal,
		"stdout":  expectedSyslog,
		"kernel":  expectedJournal,
	}

	for transport, expected := range testCases {
		// Given a message via this transport
		entry.Fields["_TRANSPORT"] = transport
		// When we format this entry
		actual := FormatEntryWithoutTime(entry)
		// Then we should get a line like the expected
		if actual != expected {
			t.Errorf("For transport '%s', Expected '%s', got '%s'",
				transport, expected, actual)
		}
	}

	// Given a message from syslog
	entry.Fields["_TRANSPORT"] = "syslog"
	// But missing the SYSLOG_PID field
	delete(entry.Fields, "SYSLOG_PID")
	// When we format this entry
	actual := FormatEntryWithoutTime(entry)
	// Then we should see the _PID field used instead
	expected := "program[42]: Hello, World"
	if actual != expected {
		t.Errorf("Ignores _PID, expected '%s', got '%s'",
			expected, actual)
	}
}
