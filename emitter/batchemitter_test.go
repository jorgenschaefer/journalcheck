package emitter

import "testing"

func TestConsume(t *testing.T) {
	// - Should call the sender after timeout

	// - Should call the sender if no timeout but maxlen entries
	// received

	// - Should call the sender only after maxduration if entries
	// are sent more often than timeout
}
