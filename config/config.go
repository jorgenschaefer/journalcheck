// Package config encapsulates journalcheck's configuration options.
// As journalcheck can be configured both using command line arguments
// and environment variables, this package provides simple accessor
// functions that return the value.
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

// FilterFile returns the name of the file containing filter rules.
// If no such option was given, ok will be false.
func FilterFile() (filename string, ok bool) {
	return optString(*flagFilterFile, "JOURNALCHECK_FILTERFILE")
}

// CursorFile returns the name of the file journalcheck should use for
// the cursor. If no such option was given, ok will be false.
func CursorFile() (string, bool) {
	return optString(*flagCursorFile, "JOURNALCHECK_CURSORFILE")
}

// RecipientAddress returns the e-mail address to be used for sending
// notification mails. If no such option was given, ok will be false.
func RecipientAddress() (string, bool) {
	return optString(*flagRecipient, "JOURNALCHECK_RECIPIENT")
}

// NumEntries returns the number of entries to show in the interactive
// use case.
func NumEntries() int {
	return *flagNumEntries
}

// OutputFormat returns the format to use for formatting journal
// entries on output.
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
