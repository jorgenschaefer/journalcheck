package config

import (
	"os"
	"strconv"
	"time"
)

func GetFilterFile() (string, bool) {
	return optionalEnvString("JOURNALCHECK_FILTERFILE")
}

func GetCursorFile() (string, bool) {
	return optionalEnvString("JOURNALCHECK_CURSORFILE")
}

func GetRecipientAddress() (string, bool) {
	return optionalEnvString("JOURNALCHECK_RECIPIENT")
}

func GetDefaultEntryCount() int {
	return envIntDefault("JOURNALCHECK_DEFAULTENTRYCOUNT", 100)
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
