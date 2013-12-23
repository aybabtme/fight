package udp_vs_tcp

import (
	"fmt"
	"net"
	"time"
)

type TCPServer struct {
}

func NewTCPServer() Server {
	return &TCPServer{}
}

func (t *TCPServer) Measure() ([]ServerMeasure, error) {
	raddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9000")
	if err != nil {
		return nil, fmt.Errorf("resolving addr, %v", err)
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
		measures = append(measures, ServerMeasure{dT, n, buf[:n], err})
		if err != nil {
			return measures, fmt.Errorf("reading, %v", err)
		}
		fmt.Printf("Read : %v\n", buf[:n])
	}

}
