package stdout

import (
	"fmt"

	"github.com/jorgenschaefer/journalcheck/event"
	"github.com/jorgenschaefer/journalcheck/notifier"
)

type StdoutNotifier int

func New() notifier.Notifier {
	return StdoutNotifier(0)
}

func (n StdoutNotifier) Receive(eventstream chan event.Event) {
	for event := range eventstream {
		fmt.Println(event.MatchString())
	}
}
