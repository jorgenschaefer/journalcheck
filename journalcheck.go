package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jorgenschaefer/journalcheck/config"
	"github.com/jorgenschaefer/journalcheck/journal"
	"github.com/jorgenschaefer/journalcheck/matcher"
)

func main() {
	var unfilteredEntries = make(chan *journal.Entry)
	var filteredEntries = make(chan *journal.Entry)
	var finalCursor = make(chan string)

	go generate(unfilteredEntries, finalCursor)
	go filter(unfilteredEntries, filteredEntries)
	recipient, ok := config.RecipientAddress()
	if ok {
		sendmail(filteredEntries, recipient)
	} else {
		writeentries(filteredEntries)
	}

	cursor, ok := <-finalCursor
	if ok {
		writeCursor(cursor)
	}
}

func generate(unfilteredEntries chan *journal.Entry, finalCursor chan string) {
	var lastEntry *journal.Entry

	j, err := journal.New()
	if err != nil {
		log.Fatal(err)
	}
	cursorfile, ok := config.CursorFile()
	if ok {
		cursor, err := ioutil.ReadFile(cursorfile)
		if err != nil {
			j.SeekLast(1)
			entry, err := j.Next()
			if err != nil {
				log.Fatal(err)
			}
			close(unfilteredEntries)
			finalCursor <- entry.Cursor
			close(finalCursor)
			fmt.Println("FOO 1")
			return
		} else {
			j.SeekCursor(string(cursor))
		}
	} else {
		j.SeekLast(config.NumEntries())
	}
	for {
		entry, err := j.Next()
		if err != nil {
			log.Fatal(err)
		}
		if entry == nil {
			close(unfilteredEntries)
			if lastEntry != nil {
				finalCursor <- lastEntry.Cursor
			}
			close(finalCursor)
			return
		}
		lastEntry = entry
		unfilteredEntries <- entry
	}
}

func filter(unfilteredEntries, filteredEntries chan *journal.Entry) {
	filterfile, ok := config.FilterFile()
	if !ok {
		for entry := range unfilteredEntries {
			filteredEntries <- entry
		}
	} else {
		m := matcher.New(filterfile)
		for entry := range unfilteredEntries {
			matches, err := m.Matches(entry.MatchString())
			if err != nil {
				log.Fatal(err)
			}
			if !matches {
				filteredEntries <- entry
			}
		}
	}
	close(filteredEntries)
}

func sendmail(filteredEntries chan *journal.Entry, recipient string) {
	count := 0
	for _ = range filteredEntries {
		count++
	}
	fmt.Printf("Sent %d entries to %s\n", count, recipient)
}

func writeentries(filteredEntries chan *journal.Entry) {
	for entry := range filteredEntries {
		switch config.OutputFormat() {
		case "short":
			fmt.Println(entry.ShortString())
		case "verbose":
			fmt.Println(entry.VerboseString())
		case "match":
			fmt.Println(entry.MatchString())
		}
	}
}

func writeCursor(cursor string) {
	cursorfile, ok := config.CursorFile()
	if ok {
		err := ioutil.WriteFile(cursorfile, ([]byte)(cursor), 0600)
		if err != nil {
			log.Fatal(err)
		}
	}
}
