package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

var flagTestMode = flag.Bool("t", false, "Test mode; emit unfiltered entries to standard output")
var flagLastEntries = flag.Int("l", 100, "The `number` of entries to parse in test mode")
var flagFilterFile = flag.String("f", "", "The filter `file` to use")

func init() {
	flag.Parse()
}

func GetFilterFile() (string, bool) {
	if *flagFilterFile != "" {
		return *flagFilterFile, true
	} else {
		return optionalEnvString("JOURNALCHECK_FILTERFILE")
	}
}

func GetCursorFile() (string, bool) {
	return optionalEnvString("JOURNALCHECK_CURSORFILE")
}

func GetRecipientAddress() (string, bool) {
	return optionalEnvString("JOURNALCHECK_RECIPIENT")
}

func GetDefaultEntryCount() int {
	return *flagLastEntries
}

func GetMaxEntriesPerBatch() int {
	return envIntDefault("JOURNALCHECK_MAXENTRIESPERBATCH", 1000)
}

func GetMaxDelayPerBatch() time.Duration {
	minutes := envIntDefault("JOURNALCHECK_MAXMINUTESPERBATCH", 60)
	return time.Duration(minutes) * time.Minute
}

func GetMaxWaitForEntries() time.Duration {
	minutes := envIntDefault("JOURNALCHECK_WAITMINUTESFORENTRIES", 60)
	return time.Duration(minutes) * time.Minute
}

func IsTestMode() bool {
	return *flagTestMode
}

func optionalEnvString(name string) (string, bool) {
	value := os.Getenv(name)
	if value == "" {
		return "", false
	} else {
		return value, true
	}
}

func envIntDefault(name string, def int) int {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	if num, err := strconv.Atoi(value); err != nil {
		return num
	} else {
		return def
	}
}
