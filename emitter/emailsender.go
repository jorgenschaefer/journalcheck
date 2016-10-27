package emitter

import (
	"log"

	"github.com/coreos/go-systemd/sdjournal"
)

type EmailSender struct {
	address string
}

func NewEmailSender(address string) *EmailSender {
	return &EmailSender{address: address}
}

func (e *EmailSender) Send(batch []*sdjournal.JournalEntry) {
	log.Printf("Sending %d entries\n", len(batch))
	// sendmail := exec.Command("/usr/sbin/sendmail", s.recipient)
	// stdin, err := sendmail.StdinPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := sendmail.Start(); err != nil {
	// 	log.Fatal(err)
	// }

	// mail := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n",
	// 	"journalcheck", s.recipient, "Journal output")
	// stdin.Write([]byte(mail))
	// for _, entry := range s.entries {
	// 	stdin.Write([]byte(formatter.Syslog(entry)))
	// 	stdin.Write([]byte("\n"))
	// }
	// stdin.Close()
	// sendmail.Wait()
}
