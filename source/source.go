package source

import "github.com/jorgenschaefer/journalcheck/event"

type Source interface {
	Emit(chan event.Event)
}
