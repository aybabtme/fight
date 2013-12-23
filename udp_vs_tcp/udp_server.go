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

func (t *UDPServer) Measure(bufferSize int) ([]ServerMeasure, error) {
	laddr, err := net.ResolveUDPAddr("udp", t.laddr)
	if err != nil {
		return nil, fmt.Errorf("resolving addr, %v", err)
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, fmt.Errorf("listening, %v", err)
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
		// don't t.gen.ValidateNext(), UDP can arrive out of order
		// and that's normal/we don't care
		measures = append(measures, ServerMeasure{dT, n, buf[:n], err})

		if err == io.EOF && !t.gen.HasNext() {
			return measures, errors.New("expected next sequence but got EOF")
		}
	}
	return measures, nil
}
