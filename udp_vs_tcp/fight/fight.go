package main

import (
	"encoding/csv"
	"fmt"
	"github.com/aybabtme/fight/udp_vs_tcp"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"time"
)

type Specification struct {
	LocalAddr  string
	RemoteAddr string
	LocalCtrl  string
	RemoteCtrl string
	Role       string
	TimeOut    int
	ByteSize   int
	CountTo    int
}

func ParseSpec() Specification {
	var spec Specification
	err := envconfig.Process("fight", &spec)
	if err != nil {
		log.Fatalf("Parsing spec, %v", err)
	}
	return spec
}

var spec = ParseSpec()

func main() {
	switch spec.Role {
	case "tcp_server":
		measureTCPServer()
	case "tcp_client":
		measureTCPClient()
	case "udp_server":
	case "udp_client":
	default:
		log.Fatalf("Invalid role, %s", spec.Role)
	}
}

func measureTCPServer() {
	lstnAddr := spec.LocalAddr
	ctrlAddr := spec.LocalCtrl
	rmtAddr := spec.RemoteAddr
	timeout := time.Millisecond * time.Duration(spec.TimeOut)
	byteSize := spec.ByteSize
	countFrom := uint64(0)
	countTo := uint64(spec.CountTo)

	filename := fmt.Sprintf("server_tcp_%s.csv",
		time.Now().Format("06_01_02_15:04:05"))

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	csvW := csv.NewWriter(f)
	defer csvW.Flush()

	gen := udp_vs_tcp.NewSequenceCounter(byteSize, countFrom, countTo)
	udp_vs_tcp.MeasureTCPServer(lstnAddr, ctrlAddr, rmtAddr, timeout, gen, csvW)
}

func measureTCPClient() {
	rmtAddr := spec.RemoteAddr
	rmtCtrlAddr := spec.RemoteCtrl
	timeout := time.Millisecond * time.Duration(spec.TimeOut)
	byteSize := spec.ByteSize
	countFrom := uint64(0)
	countTo := uint64(spec.CountTo)

	filename := fmt.Sprintf("client_tcp_%s.csv",
		time.Now().Format("06_01_02_15:04:05"))

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	csvW := csv.NewWriter(f)
	defer csvW.Flush()

	gen := udp_vs_tcp.NewSequenceCounter(byteSize, countFrom, countTo)
	udp_vs_tcp.MeasureTCPClient(rmtAddr, rmtCtrlAddr, timeout, gen, csvW)
}
