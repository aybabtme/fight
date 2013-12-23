package udp_vs_tcp

import (
	"time"
)

type ServerMeasure struct {
	dT   time.Duration
	size int
	msg  []byte
	err  error
}

type Server interface {
	Measure() ([]ServerMeasure, error)
}
