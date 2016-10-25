package main

import (
	"io/ioutil"

	"github.com/coreos/go-systemd/sdjournal"

	"github.com/jorgenschaefer/journalcheck/config"
	"github.com/jorgenschaefer/journalcheck/emitter"
	"github.com/jorgenschaefer/journalcheck/filter"
	"github.com/jorgenschaefer/journalcheck/journal"
)

func main() {
	producer := getProducer()
	consumer := getConsumer()

	entries := make(chan *sdjournal.JournalEntry)
	go producer.Produce(entries)
	consumer.Consume(entries)
}

func getProducer() *journal.Producer {
	p := journal.NewProducer()
	if config.IsTestMode() {
		p.SeekLast(uint64(config.DefaultEntryCount()))
		p.Terminate = true
		return p
	} else if cursorfile, ok := config.CursorFile(); ok {
		if cursor, err := ioutil.ReadFile(cursorfile); err == nil {
			p.SeekCursor(string(cursor))
			return p
		}
	}
	p.SeekLast(0)
	return p
}

func getConsumer() emitter.Emitter {
	var e emitter.Emitter
	if address, ok := config.RecipientAddress(); ok {
		e = getEmailEmitter(address)
	} else {
		e = emitter.NewStdoutEmitter()
	}
	if filename, ok := config.FilterFile(); ok {
		e.SetFilter(filter.NewRegexpFilter(filename))
	}
	if filename, ok := config.CursorFile(); ok {
		e.SetCursorFile(filename)
	}
	return e
}

func getEmailEmitter(address string) emitter.Emitter {
	e := emitter.NewEmailEmitter(address)
	e.MaxLen = config.GetMaxEntriesPerBatch()
	e.MaxDuration = config.GetMaxWaitForEntries()
	return e
}
