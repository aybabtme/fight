package udp_vs_tcp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"time"
)

type ServerMeasure struct {
	t    time.Time
	size int
	msg  []byte
	err  error
}

type Server interface {
	Measure(bufferSize int, kill chan struct{}) (chan *ServerMeasure, error)
}

type ControlSrv struct {
	laddr string
	raddr string
	k     chan struct{}
	l     net.Listener
}

func NewControlServer(ctrlAddr, rmtAddr string, kill chan struct{}) *ControlSrv {
	return &ControlSrv{laddr: ctrlAddr, raddr: rmtAddr, k: kill}
}

func (c *ControlSrv) Start() error {
	l, err := net.Listen("tcp", c.laddr)
	if err != nil {
		return fmt.Errorf("creating listener for http, %v", err)
	}
	c.l = l

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", c.ping)
	mux.HandleFunc("/kill", c.kill)

	return http.Serve(c.l, mux)
}

func (c *ControlSrv) Close() error {
	defer c.l.Close()
	c.k <- struct{}{}
	resp, err := http.Get(path.Join(c.raddr, "/kill"))
	if err != nil {
		return fmt.Errorf("sending kill message to other, %v", err)
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response to kill request, %v", err)
	}

	want := []byte("ok\n")
	if !bytes.Equal(data, want) {
		return fmt.Errorf("want '%s' but was '%s'", want, data)
	}
	return nil
}

func (c *ControlSrv) ping(rw http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintln(rw, "pong")
	if err != nil {
		log.Fatalf("Responding to `/ping`, %v", err)
	}
}

func (c *ControlSrv) kill(rw http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintln(rw, "ok")
	if err != nil {
		log.Fatalf("Responding to `/kill`, %v", err)
	}

	if err := c.Close(); err != nil {
		log.Fatalf("Closing down after `/kill`, %v", err)
	}
}
