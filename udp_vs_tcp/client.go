package udp_vs_tcp

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"time"
)

type ClientMeasure struct {
	t    time.Time
	dT   time.Duration
	size int
	err  error
}

type Client struct {
	net         string
	rmt         net.Addr
	rmtCtrlAddr string
	timeout     time.Duration
	generator   Generator
}

func NewClient(network string, addr net.Addr, rmtCtrlAddr string, timeout time.Duration, gen Generator) *Client {
	return &Client{
		net:         network,
		rmt:         addr,
		rmtCtrlAddr: rmtCtrlAddr,
		timeout:     timeout,
		generator:   gen,
	}
}

func (c *Client) Send(kill chan struct{}) (chan *ClientMeasure, error) {
	conn, err := net.DialTimeout(c.net, c.rmt.String(), c.timeout)
	if err != nil {
		return nil, fmt.Errorf("dialing, %v", err)
	}

	measures := make(chan *ClientMeasure)

	go func(mChan chan *ClientMeasure) {
		defer conn.Close()
		defer close(mChan)

		var last time.Time
		for c.generator.HasNext() {
			select {
			case <-kill:
				log.Printf("Client received kill request")
				return
			default:
				fmt.Printf("%d/%d\r", c.generator.Done(), c.generator.Total())
				last = time.Now()
				m, err := c.write(conn, last)
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

		log.Printf("Done sending, %d/%d", c.generator.Done(), c.generator.Total())

		rmtCtrl := "http://" + path.Join(c.rmtCtrlAddr, "/kill")
		log.Printf("Done sending, GET %v", rmtCtrl)
		_, err := http.Get(rmtCtrl)
		if err != nil {
			log.Fatalf("Sending GET /kill to server, %v", err)
		}
	}(measures)

	return measures, nil
}

func (c *Client) write(conn net.Conn, last time.Time) (*ClientMeasure, error) {
	n, err := conn.Write(c.generator.Next())
	now := time.Now()
	return &ClientMeasure{now, now.Sub(last), n, err}, err
}
