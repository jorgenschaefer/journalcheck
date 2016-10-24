package source

import (
	"log"

	"github.com/coreos/go-systemd/sdjournal"
)

type Producer struct {
	journal    *sdjournal.Journal
	cursorFile *string
	rewind     uint64
}

func NewProducer() *Producer {
	producer := new(Producer)
	journal, err := sdjournal.NewJournal()
	if err != nil {
		log.Fatal(err)
	}
	producer.journal = journal
	return producer
}

func (p *Producer) SeekCursor(cursor string) {
	if err := p.journal.SeekCursor(cursor); err != nil {
		log.Fatal(err)
	}
	// Move to the position of the cursor, else we'd see the last
	// entry again
	_, err := p.journal.Next()
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Producer) SeekLast(rewind uint64) {
	if err := p.journal.SeekTail(); err != nil {
		log.Fatal(err)
	}
	_, err := p.journal.PreviousSkip(rewind)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Producer) Produce(entries chan *sdjournal.JournalEntry) {
	for {
		n, err := p.journal.Next()
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			p.journal.Wait(sdjournal.IndefiniteWait)
			continue
		}
		entry, err := p.journal.GetEntry()
		if err != nil {
			log.Fatal(err)
		}
		entries <- entry
	}
}

func next(journal *sdjournal.Journal) *sdjournal.JournalEntry {
	n, err := journal.Next()
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		return nil
	}
	entry, err := journal.GetEntry()
	if err != nil {
		log.Fatal(err)
	}
	return entry
}
