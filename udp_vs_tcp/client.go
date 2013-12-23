package udp_vs_tcp

import (
	"fmt"
	"net"
	"time"
)

type ClientMeasure struct {
	dT   time.Duration
	size int
	err  error
}

type Client struct {
	net       string
	dest      net.Addr
	timeout   time.Duration
	generator Generator
}

func NewClient(network string, addr net.Addr, timeout time.Duration, gen Generator) *Client {
	return &Client{
		net:       network,
		dest:      addr,
		timeout:   timeout,
		generator: gen,
	}
}

func (t *Client) Send(size int) (measures []ClientMeasure, closeErr error) {
	conn, err := net.DialTimeout(t.net, t.dest.String(), t.timeout)
	if err != nil {
		return nil, fmt.Errorf("dialing, %v", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
		if closeErr != nil {
			closeErr = fmt.Errorf("%v, prior error : %v", err, closeErr)
		} else {
			closeErr = err
		}
	}()

	measures = make([]ClientMeasure, size)

	for i := range measures {
		start := time.Now()
		n, err := conn.Write(t.generator.Next())
		dT := time.Since(start)
		measures[i] = ClientMeasure{dT, n, err}

		if err != nil {
			opErr, ok := err.(*net.OpError)
			if ok && opErr.Temporary() {
				continue
			}
			return measures[:i+1], err
		}
	}
	return
}
