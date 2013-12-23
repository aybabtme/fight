package udp_vs_tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TCPServer struct {
	raddr string
	gen   Generator
}

func NewTCPServer(raddr string, gen Generator) Server {
	return &TCPServer{raddr, gen}
}

func (t *TCPServer) Measure() ([]ServerMeasure, error) {
	raddr, err := net.ResolveTCPAddr("tcp", t.raddr)
	if err != nil {
		return nil, fmt.Errorf("reso.lving addr, %v", err)
	}

	l, err := net.ListenTCP("tcp", raddr)
	if err != nil {
		return nil, fmt.Errorf("listening, %v", err)
	}
	defer l.Close()

	conn, err := l.AcceptTCP()
	if err != nil {
		return nil, fmt.Errorf("accepting, %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 1<<14)
	measures := make([]ServerMeasure, 0)
	for {
		start := time.Now()
		n, err := conn.Read(buf)
		dT := time.Since(start)
		if err != nil && err != io.EOF {
			measures = append(measures, ServerMeasure{dT, n, buf[:n], err})
			return measures, fmt.Errorf("reading, %v", err)
		}
		if !t.gen.ValidateNext(buf[:n]) {
			err = errors.New("failed to read expected sequence")
		}
		measures = append(measures, ServerMeasure{dT, n, buf[:n], err})

		if err == io.EOF && !t.gen.HasNext() {
			return measures, errors.New("expected next sequence but got EOF")
		} else if err != io.EOF && t.gen.HasNext() {
			return measures, errors.New("expected EOF but got nothing")
		}
	}
	return measures, nil
}
