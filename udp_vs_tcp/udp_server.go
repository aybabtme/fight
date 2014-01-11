package udp_vs_tcp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type UDPServer struct {
	laddr   string
	gen     Generator
	timeout time.Duration
}

func NewUDPServer(laddr string, timeout time.Duration, gen Generator) Server {
	return &UDPServer{laddr, gen, timeout}
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

	// Kill routine
	go func() {
		<-kill
		log.Printf("Server received kill request")
		err := conn.Close()
		if err != nil {
			log.Fatalf("Closing connection, %v", err)
		}
	}()

	// Read routine
	measures := make(chan *ServerMeasure)
	go func(mChan chan *ServerMeasure) {
		defer conn.Close()
		defer close(mChan)

		buf := make([]byte, bufferSize)
		var last time.Time
		for {
			last = time.Now()
			m, err := u.read(conn, buf, last)
			if err != nil {
				panic(err)
			}
			mChan <- m
		}
	}(measures)
	return measures, nil
}

func (u *UDPServer) read(conn *net.UDPConn, buf []byte, last time.Time) (*ServerMeasure, error) {
	err := conn.SetDeadline(time.Now().Add(u.timeout))
	if err != nil {
		return nil, fmt.Errorf("setting deadline for next read, %v", err)
	}

	n, err := conn.Read(buf)
	now := time.Now()

	m := &ServerMeasure{now, now.Sub(last), n, buf[:n], err}

	if err != nil && err != io.EOF {
		return m, fmt.Errorf("reading, %v", err)
	}

	if err == io.EOF && !u.gen.HasNext() {
		return m, errors.New("expected next sequence but got EOF")
	}

	return m, nil
}
