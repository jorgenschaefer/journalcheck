package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/quotedprintable"
	"os"
	"os/exec"
	"time"

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
		m, err := matcher.New(filterfile)
		if err != nil {
			log.Fatal(err)
		}
		for entry := range unfilteredEntries {
			matches := m.Matches(entry.MatchString())
			if !matches {
				filteredEntries <- entry
			}
		}
	}
	close(filteredEntries)
}

func sendmail(filteredEntries chan *journal.Entry, recipient string) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	cmd := exec.Command("/usr/sbin/sendmail", recipient)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	mail := quotedprintable.NewWriter(stdin)
	_, err = fmt.Fprintf(mail, `From: logcheck@%s
To: %s
Subject: Journalcheck at %s
Content-Type: text/plain
Content-Transfer-Encoding: quoted-printable
MIME: 1.0

This email is sent by journalcheck. If you no longer wish to receive
such mail, you can either deinstall the journalcheck package or modify
the configuration.

`, hostname, recipient, time.Now().Format(time.RFC822))
	if err != nil {
		log.Fatal(err)
	}

	for entry := range filteredEntries {
		fmt.Fprintln(mail, entry.ShortString())
	}

	if err := stdin.Close(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
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
