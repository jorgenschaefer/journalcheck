package notifier

import "github.com/jorgenschaefer/journalcheck/event"

type Notifier interface {
	Receive(chan event.Event)
}
