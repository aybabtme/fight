package udp_vs_tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TCPServer struct {
	laddr string
	gen   Generator
}

func NewTCPServer(laddr string, gen Generator) Server {
	return &TCPServer{laddr, gen}
}

func (t *TCPServer) Measure(bufferSize int) ([]ServerMeasure, error) {
	laddr, err := net.ResolveTCPAddr("tcp", t.laddr)
	if err != nil {
		return nil, fmt.Errorf("resolving addr, %v", err)
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, fmt.Errorf("listening, %v", err)
	}
	defer l.Close()

	conn, err := l.AcceptTCP()
	if err != nil {
		return nil, fmt.Errorf("accepting, %v", err)
	}
	defer conn.Close()

	buf := make([]byte, bufferSize)
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
