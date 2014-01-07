package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/aybabtme/fight/udp_vs_tcp"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Specification struct {
	LocalAddr  string `json:"localAddress"`
	RemoteAddr string `json:"remoteAddress"`
	LocalCtrl  string `json:"localControlAddress"`
	RemoteCtrl string `json:"remoteControlAddress"`
	Role       string `json:"localRole"`
	TimeOut    int    `json:"timeout"`
	ByteSize   int    `json:"bytesize"`
	CountTo    int    `json:"countTo"`
}

var spec = ParseSpec()

func ParseSpec() Specification {
	f, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return WriteSpec()
		}
		panic(err)
	}
	defer f.Close()
	d := json.NewDecoder(f)

	var spec Specification

	if err := d.Decode(&spec); err != nil {
		panic(err)
	}
	return spec
}

func WriteSpec() Specification {
	spec := Specification{
		LocalAddr:  "127.0.0.1:6060",
		RemoteAddr: "127.0.0.1:6060",
		LocalCtrl:  "127.0.0.1:8080",
		RemoteCtrl: "127.0.0.1:8080",
		Role:       "tcp_server",
		TimeOut:    15000,
		ByteSize:   8,
		CountTo:    10000,
	}

	data, err := json.MarshalIndent(&spec, "", "   ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("config.json", data, 0660)
	if err != nil {
		panic(err)
	}

	return spec
}

func main() {
	switch spec.Role {
	case "tcp_server":
		measureTCPServer()
	case "tcp_client":
		measureTCPClient()
	case "udp_server":
		measureUDPServer()
	case "udp_client":
		measureUDPClient()
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
		time.Now().Format("2006_01_02_15:04:05"))

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
		time.Now().Format("2006_01_02_15:04:05"))

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

func measureUDPServer() {
	lstnAddr := spec.LocalAddr
	ctrlAddr := spec.LocalCtrl
	rmtAddr := spec.RemoteAddr
	timeout := time.Millisecond * time.Duration(spec.TimeOut)
	byteSize := spec.ByteSize
	countFrom := uint64(0)
	countTo := uint64(spec.CountTo)

	filename := fmt.Sprintf("server_udp_%s.csv",
		time.Now().Format("2006_01_02_15:04:05"))

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	csvW := csv.NewWriter(f)
	defer csvW.Flush()

	gen := udp_vs_tcp.NewSequenceCounter(byteSize, countFrom, countTo)
	udp_vs_tcp.MeasureUDPServer(lstnAddr, ctrlAddr, rmtAddr, timeout, gen, csvW)
}

func measureUDPClient() {
	rmtAddr := spec.RemoteAddr
	rmtCtrlAddr := spec.RemoteCtrl
	timeout := time.Millisecond * time.Duration(spec.TimeOut)
	byteSize := spec.ByteSize
	countFrom := uint64(0)
	countTo := uint64(spec.CountTo)

	filename := fmt.Sprintf("client_udp_%s.csv",
		time.Now().Format("2006_01_02_15:04:05"))

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	csvW := csv.NewWriter(f)
	defer csvW.Flush()

	gen := udp_vs_tcp.NewSequenceCounter(byteSize, countFrom, countTo)
	udp_vs_tcp.MeasureUDPClient(rmtAddr, rmtCtrlAddr, timeout, gen, csvW)
}
