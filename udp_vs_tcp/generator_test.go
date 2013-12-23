package udp_vs_tcp

import (
	"testing"
)

func TestSequenceGenerator(t *testing.T) {
	f := NewSequenceCounter(1<<14, 0, 1<<16)
	g := NewSequenceCounter(1<<14, 0, 1<<16)

	var next []byte
	for g.HasNext() && f.HasNext() {
		next = g.Next()
		if !f.ValidateNext(next) {
			t.Fatalf("Mismatch at %d/%d: % x", f.Done(), f.Total(), next)
		}
	}
}
