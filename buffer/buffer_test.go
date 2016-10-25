package buffer

import (
	"testing"

	"github.com/coreos/go-systemd/sdjournal"
)

func TestIsEmpty(t *testing.T) {
	if !NewBuffer(23).IsEmpty() {
		t.Errorf("New buffer should be empty, but isn't")
	}

	buffer, _ := NewBuffer(23).Append(nil)
	if buffer.IsEmpty() {
		t.Errorf("Buffer with one element should not be empty, but is")
	}
}

func TestAppend(t *testing.T) {
	v := &sdjournal.JournalEntry{}
	b, err := NewBuffer(1).Append(v)
	if err != nil {
		t.Errorf("A one-element buffer should accept an element")
	}
	if b[0] != v {
		t.Errorf("The first element should be %v, not %v", v, b[0])
	}
	_, err = b.Append(nil)
	if err == nil {
		t.Errorf("A full buffer could be appended to")
	}

}

func TestIsFull(t *testing.T) {
	b := NewBuffer(1)
	if b.IsFull() {
		t.Errorf("A new buffer should not be full")
	}
	b, _ = b.Append(nil)
	if !b.IsFull() {
		t.Errorf("A one-element buffer with one element is not full")
	}
}

func TestClear(t *testing.T) {
	b := NewBuffer(1)
	b, _ = b.Append(nil)
	b = b.Clear()
	if !b.IsEmpty() {
		t.Errorf("A cleared buffer is not empty")
	}
}
