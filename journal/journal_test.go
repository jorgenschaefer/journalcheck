package journal

import "testing"

func TestMatchString(t *testing.T) {
	// Given an entry with a systemd unit
	fields := map[string]string{
		"_SYSTEMD_UNIT":     "test.service",
		"_COMM":             "/bin/test",
		"SYSLOG_IDENTIFIER": "test",
		"MESSAGE":           "This is a test",
	}
	e := &Entry{Fields: fields}
	// Then the match string should include the systemd unit
	expectEqual(t, "test.service: This is a test", e.MatchString())
	// Given an entry without a systemd unit, but a command
	delete(fields, "_SYSTEMD_UNIT")
	// Then the match string should include the command
	expectEqual(t, "/bin/test: This is a test", e.MatchString())
	// Given an entry without a command but a syslog identifier
	delete(fields, "_COMM")
	// Then the match string should include the syslog identifier
	expectEqual(t, "test: This is a test", e.MatchString())
}

func TestShortString(t *testing.T) {
	// Given an entry with a systemd unit
	fields := map[string]string{
		"_SYSTEMD_UNIT":     "test.service",
		"_COMM":             "/bin/test",
		"SYSLOG_IDENTIFIER": "test",
		"MESSAGE":           "This is a test",
		"_HOSTNAME":         "testhost",
		"SYSLOG_PID":        "12345",
		"_PID":              "54321",
	}
	e := &Entry{Fields: fields, RealtimeTimestamp: 1477764000 * 1000 * 1000}
	// Then the match string should include the systemd unit
	expectEqual(t, "Oct 29 20:00:00 testhost test.service[12345]: This is a test", e.ShortString())
	// Given an entry without a syslog PID
	delete(fields, "SYSLOG_PID")
	// Then the entry should use _PID
	expectEqual(t, "Oct 29 20:00:00 testhost test.service[54321]: This is a test", e.ShortString())
	// Given an entry without a systemd unit, but a command
	delete(fields, "_SYSTEMD_UNIT")
	// Then the match string should include the command
	expectEqual(t, "Oct 29 20:00:00 testhost /bin/test[54321]: This is a test", e.ShortString())
	// Given an entry without a command but a syslog identifier
	delete(fields, "_COMM")
	// Then the match string should include the syslog identifier
	expectEqual(t, "Oct 29 20:00:00 testhost test[54321]: This is a test", e.ShortString())
}

func expectEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected the match string to be '%v', but it was '%v'", expected, actual)
	}
}
