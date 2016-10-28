package config

import (
	"flag"
	"log"
	"os"
)

var flagCursorFile = flag.String("c", "", "A `file` to store the last read cursor in")
var flagNumEntries = flag.Int("n", 100, "The `number` of entries to parse in test mode")
var flagFilterFile = flag.String("f", "", "The filter `file` to use")
var flagOutputFormat = flag.String("o", "short", "Output format (one of: short, verbose, match)")
var flagRecipient = flag.String("r", "", "Send e-mails to this e-mail `address`")

func init() {
	flag.Parse()
	switch *flagOutputFormat {
	case "short":
	case "verbose":
	case "match":
		break
	default:
		log.Fatalf("Bad value for output format (-o): %s. Should be one of short, verbose or match", *flagOutputFormat)
	}
}

func FilterFile() (string, bool) {
	return optString(*flagFilterFile, "JOURNALCHECK_FILTERFILE")
}

func CursorFile() (string, bool) {
	return optString(*flagCursorFile, "JOURNALCHECK_CURSORFILE")
}

func RecipientAddress() (string, bool) {
	return optString(*flagRecipient, "JOURNALCHECK_RECIPIENT")
}

func NumEntries() int {
	return *flagNumEntries
}

func OutputFormat() string {
	return *flagOutputFormat
}

func optString(option, envvar string) (string, bool) {
	if option != "" {
		return option, true
	}
	value := os.Getenv(envvar)
	if value != "" {
		return value, true
	}
	return "", false
}
