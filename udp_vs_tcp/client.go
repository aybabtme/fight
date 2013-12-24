package udp_vs_tcp

import (
	"fmt"
	"net"
	"time"
)

type ClientMeasure struct {
	t    time.Time
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

func (c *Client) Send(size int, kill chan struct{}) (chan *ClientMeasure, error) {
	conn, err := net.DialTimeout(c.net, c.dest.String(), c.timeout)
	if err != nil {
		return nil, fmt.Errorf("dialing, %v", err)
	}

	measures := make(chan *ClientMeasure, size)

	go func(mChan chan *ClientMeasure) {
		defer conn.Close()
		defer close(mChan)

		for c.generator.HasNext() {
			select {
			case <-kill:
				return
			default:
				m, err := c.write(conn)
				mChan <- m
				if err == nil {
					continue
				}

				opErr, ok := err.(*net.OpError)
				if !ok || !opErr.Temporary() {
					panic(err)
				}
			}

		}
	}(measures)

	return measures, nil
}

func (c *Client) write(conn net.Conn) (*ClientMeasure, error) {
	n, err := conn.Write(c.generator.Next())
	return ClientMeasure{time.Now(), n, err}, err
}
