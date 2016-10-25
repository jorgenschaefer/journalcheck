package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

var flagTestMode = flag.Bool("t", false, "Test mode; emit unfiltered entries to standard output")
var flagLastEntries = flag.Int("l", 100, "The `number` of entries to parse in test mode")
var flagFilterFile = flag.String("f", "", "The filter `file` to use")

func init() {
	flag.Parse()
	_, ok := RecipientAddress()
	if !IsTestMode() && !ok {
		log.Fatal("Please either specify test mode (-t) or provide a recipient address")
	}
}

func FilterFile() (string, bool) {
	if *flagFilterFile != "" {
		return *flagFilterFile, true
	} else {
		return optionalEnvString("JOURNALCHECK_FILTERFILE")
	}
}

func CursorFile() (string, bool) {
	return optionalEnvString("JOURNALCHECK_CURSORFILE")
}

func RecipientAddress() (string, bool) {
	return optionalEnvString("JOURNALCHECK_RECIPIENT")
}

func DefaultEntryCount() int {
	return *flagLastEntries
}

func MaxEntriesPerBatch() int {
	return envIntDefault("JOURNALCHECK_MAXENTRIESPERBATCH", 1000)
}

func MaxDelayPerBatch() time.Duration {
	minutes := envIntDefault("JOURNALCHECK_MAXMINUTESPERBATCH", 60)
	return time.Duration(minutes) * time.Minute
}

func MaxWaitForEntries() time.Duration {
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
