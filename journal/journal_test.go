package journal

import (
	"testing"
	"time"
)

func TestMatchString(t *testing.T) {
	// Given an entry with a syslog identifier
	fields := map[string]string{
		"_SYSTEMD_UNIT":     "test.service",
		"_COMM":             "/bin/test",
		"SYSLOG_IDENTIFIER": "test",
		"MESSAGE":           "This is a test",
	}
	e := &Entry{Fields: fields}
	// Then the match string should include the syslog identifier
	expectEqual(t, "test: This is a test", e.MatchString())
	// Given an entry without the syslog identifier
	delete(fields, "SYSLOG_IDENTIFIER")
	// Then the match string should include the systemd unit
	expectEqual(t, "test.service: This is a test", e.MatchString())
	// Given an entry without a systemd unit, but a command
	delete(fields, "_SYSTEMD_UNIT")
	// Then the match string should include the command
	expectEqual(t, "/bin/test: This is a test", e.MatchString())
	// Given an entry without a command but a syslog identifier
	delete(fields, "_COMM")
}

func TestShortString(t *testing.T) {
	// Given an entry with a syslog identifier
	fields := map[string]string{
		"_SYSTEMD_UNIT":     "test.service",
		"_COMM":             "/bin/test",
		"SYSLOG_IDENTIFIER": "test",
		"MESSAGE":           "This is a test",
		"_HOSTNAME":         "testhost",
		"SYSLOG_PID":        "12345",
		"_PID":              "54321",
	}
	ts := time.Date(2016, 10, 29, 20, 0, 0, 0, time.Local).UnixNano() / 1000
	e := &Entry{Fields: fields, RealtimeTimestamp: uint64(ts)}
	// Then the short string should include the syslog identifier
	expectEqual(t, "Oct 29 20:00:00 testhost test[12345]: This is a test", e.ShortString())

	// Given an entry without the syslog identifier
	delete(fields, "SYSLOG_IDENTIFIER")
	// Then the short string should include the systemd unit
	expectEqual(t, "Oct 29 20:00:00 testhost test.service[12345]: This is a test", e.ShortString())

	// Given an entry without a syslog PID
	delete(fields, "SYSLOG_PID")
	// Then the short string should use _PID
	expectEqual(t, "Oct 29 20:00:00 testhost test.service[54321]: This is a test", e.ShortString())

	// Given an entry without a systemd unit, but a command
	delete(fields, "_SYSTEMD_UNIT")
	// Then the short string should include the command
	expectEqual(t, "Oct 29 20:00:00 testhost /bin/test[54321]: This is a test", e.ShortString())

	// Given an entry with newlines
	fields["MESSAGE"] = "This is a test\nwith\nnewlines"
	// Then the short string should use spaces instead of the
	// newline
	expectEqual(t, "Oct 29 20:00:00 testhost /bin/test[54321]: This is a test with newlines", e.ShortString())
}

func expectEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected the match string to be '%v', but it was '%v'", expected, actual)
	}
}
