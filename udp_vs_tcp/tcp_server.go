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

func (t *TCPServer) Measure(bufferSize int, kill chan struct{}) (chan *ServerMeasure, error) {
	laddr, err := net.ResolveTCPAddr("tcp", t.laddr)
	if err != nil {
		return nil, fmt.Errorf("resolving addr, %v", err)
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, fmt.Errorf("listening, %v", err)
	}

	conn, err := l.AcceptTCP()
	if err != nil {
		return nil, fmt.Errorf("accepting, %v", err)
	}

	measures := make(chan *ServerMeasure)

	go func(mChan chan *ServerMeasure) {
		defer l.Close()
		defer conn.Close()
		defer close(mChan)

		buf := make([]byte, bufferSize)

		for {
			select {
			case <-kill:
				return
			default:
				m, err := t.read(conn, buf)
				if err != nil {
					panic(err)
				}
				mChan <- m
			}
		}
	}(measures)

	return measures, nil
}
func (t *TCPServer) read(conn *net.TCPConn, buf []byte) (*ServerMeasure, error) {
	n, err := conn.Read(buf)
	now := time.Now()

	m := &ServerMeasure{now, n, buf[:n], err}

	if err != nil && err != io.EOF {
		return m, fmt.Errorf("reading, %v", err)
	}

	if !t.gen.ValidateNext(buf[:n]) {
		return m, errors.New("failed to read expected sequence")
	}

	if err == io.EOF && !t.gen.HasNext() {
		return m, errors.New("expected next sequence but got EOF")
	}

	if err != io.EOF && t.gen.HasNext() {
		return m, errors.New("expected EOF but got nothing")
	}

	return m, nil
}
