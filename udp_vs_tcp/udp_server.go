package udp_vs_tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type UDPServer struct {
	laddr string
	gen   Generator
}

func NewUDPServer(laddr string, gen Generator) Server {
	return &UDPServer{laddr, gen}
}

func (u *UDPServer) Measure(bufferSize int, kill chan struct{}) (chan *ServerMeasure, error) {
	laddr, err := net.ResolveUDPAddr("udp", u.laddr)
	if err != nil {
		return nil, fmt.Errorf("resolving addr, %v", err)
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, fmt.Errorf("listening, %v", err)
	}

	measures := make(chan *ServerMeasure)

	go func(mChan chan *ServerMeasure) {
		defer conn.Close()
		defer close(mChan)

		buf := make([]byte, bufferSize)
		for {
			select {
			case <-kill:
				return
			default:
				m, err := u.read(conn, buf)
				if err != nil {
					panic(err)
				}
				mChan <- m
			}
		}
	}(measures)
	return measures, nil
}

func (u *UDPServer) read(conn *net.UDPConn, buf []byte) (*ServerMeasure, error) {
	n, err := conn.Read(buf)
	now := time.Now()

	m := &ServerMeasure{now, n, buf[:n], err}

	if err != nil && err != io.EOF {
		return m, fmt.Errorf("reading, %v", err)
	}

	if err == io.EOF && !u.gen.HasNext() {
		return m, errors.New("expected next sequence but got EOF")
	}

	return m, nil
}
