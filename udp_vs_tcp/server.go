package udp_vs_tcp

import (
	"time"
)

type ServerMeasure struct {
	t    time.Time
	size int
	msg  []byte
	err  error
}

type Server interface {
	Measure(bufferSize int) ([]ServerMeasure, error)
}
