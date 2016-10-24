package main

import (
	"io/ioutil"

	"github.com/coreos/go-systemd/sdjournal"

	"github.com/jorgenschaefer/journalcheck/config"
	"github.com/jorgenschaefer/journalcheck/emitter"
	"github.com/jorgenschaefer/journalcheck/filter"
	"github.com/jorgenschaefer/journalcheck/source"
)

func main() {
	producer := getProducer()
	consumer := getConsumer()

	entries := make(chan *sdjournal.JournalEntry)
	go producer.Produce(entries)
	consumer.Consume(entries)
}

func getProducer() *source.Producer {
	var p *source.Producer
	p = source.NewProducer()
	if cursorfile, ok := config.GetCursorFile(); ok {
		cursor, err := ioutil.ReadFile(cursorfile)
		if err == nil {
			p.SeekCursor(string(cursor))
			return p
		}
	}
	p.SeekLast(uint64(config.GetDefaultEntryCount()))
	return p
}

func getConsumer() emitter.Emitter {
	var e emitter.Emitter
	if address, ok := config.GetRecipientAddress(); ok {
		e = getEmailEmitter(address)
	} else {
		e = emitter.NewStdoutEmitter()
	}
	if filename, ok := config.GetFilterFile(); ok {
		e.SetFilter(filter.NewRegexpFilter(filename))
	}
	if filename, ok := config.GetCursorFile(); ok {
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
