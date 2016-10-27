package main

import (
	"log"

	"github.com/jorgenschaefer/journalcheck/config"
	"github.com/jorgenschaefer/journalcheck/event"
	"github.com/jorgenschaefer/journalcheck/notifier"
	"github.com/jorgenschaefer/journalcheck/notifier/stdout"
	"github.com/jorgenschaefer/journalcheck/source"
	"github.com/jorgenschaefer/journalcheck/source/journal"
)

func main() {
	source := getSource()
	notifier := getNotifier()

	eventstream := make(chan event.Event)
	go source.Emit(eventstream)
	notifier.Receive(eventstream)
}

func getSource() source.Source {
	filterfile, ok := config.FilterFile()
	var ff *string
	if ok {
		ff = &filterfile
	} else {
		ff = nil
	}
	if config.IsTestMode() {
		entryCount := config.DefaultEntryCount()
		return journal.NewDelimitedSource(entryCount, ff)
	} else {
		cursorfile, ok := config.CursorFile()
		if !ok {
			log.Fatal("Please either specify test mode or provide a cursor file name")
		}
		return journal.NewCursorSource(cursorfile, ff)
	}
}

func getNotifier() notifier.Notifier {
	if _, ok := config.RecipientAddress(); ok {
		panic("E-mail notifier not implemented yet")
	} else {
		return stdout.New()
	}
}
