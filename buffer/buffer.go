package buffer

import (
	"errors"

	"github.com/coreos/go-systemd/sdjournal"
)

type buffer []*sdjournal.JournalEntry

func NewBuffer(maxLen int) buffer {
	return make(buffer, 0, maxLen)
}

func (b buffer) IsEmpty() bool {
	return len(b) == 0
}

func (b buffer) IsFull() bool {
	return len(b) == cap(b)
}

func (b buffer) Append(e *sdjournal.JournalEntry) (buffer, error) {
	if cap(b) == len(b) {
		return nil, errors.New("Buffer is full")
	} else {
		return append(b, e), nil
	}
}

func (b buffer) Clear() buffer {
	for i := range b {
		b[i] = nil
	}
	return b[0:0]
}
