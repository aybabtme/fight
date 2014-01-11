package udp_vs_tcp

import (
	"encoding/csv"
	"log"
	"net"
	"strconv"
	"time"
)

func MeasureTCPServer(lstnAddr, ctrlAddr, rmtAddr string, timeout time.Duration, gen Generator, csvW *csv.Writer) {

	kill := make(chan struct{})

	log.Printf("Starting control server")
	ctrl := NewControlServer(ctrlAddr, rmtAddr, kill)
	go func() {
		err := ctrl.Start()
		if err != nil {
			log.Fatalf("Starting control server, %v", err)
		}
	}()
	defer ctrl.Close()

	log.Printf("Starting TCP server with timeout %v", timeout)
	srv := NewTCPServer(lstnAddr, timeout, gen)

	log.Print("Starting to listen")
	measures, err := srv.Measure(gen.Size(), kill)
	if err != nil {
		log.Fatalf("Starting measurement server (closeErr=%v), %v", ctrl.Close(), err)
	}

	log.Printf("Collecting measurements")
	csvW.Write([]string{
		"Time",
		"dT",
		"MessageSize",
		"Message",
		"Error",
	})
	for measure := range measures {
		var errStr string
		if measure.err != nil {
			errStr = measure.err.Error()
		}
		csvW.Write([]string{
			measure.t.String(),
			strconv.FormatInt(measure.dT.Nanoseconds(), 10),
			strconv.Itoa(measure.size),
			string(measure.msg),
			errStr,
		})
	}
}

func MeasureUDPServer(lstnAddr, ctrlAddr, rmtAddr string, timeout time.Duration, gen Generator, csvW *csv.Writer) {

	kill := make(chan struct{})

	log.Printf("Starting control server")
	ctrl := NewControlServer(ctrlAddr, rmtAddr, kill)
	go func() {
		err := ctrl.Start()
		if err != nil {
			log.Fatalf("Starting control server, %v", err)
		}
	}()
	defer ctrl.Close()

	log.Printf("Starting UDP server with timeout %v", timeout)
	srv := NewUDPServer(lstnAddr, timeout, gen)

	log.Print("Starting to listen")
	measures, err := srv.Measure(gen.Size(), kill)
	if err != nil {
		log.Fatalf("Starting measurement server (closeErr=%v), %v", ctrl.Close(), err)
	}

	log.Printf("Collecting measurements")
	csvW.Write([]string{
		"Time",
		"dT",
		"MessageSize",
		"Message",
		"Error",
	})
	for measure := range measures {
		var errStr string
		if measure.err != nil {
			errStr = measure.err.Error()
		}
		csvW.Write([]string{
			measure.t.String(),
			strconv.FormatInt(measure.dT.Nanoseconds(), 10),
			strconv.Itoa(measure.size),
			string(measure.msg),
			errStr,
		})
	}
}

func MeasureTCPClient(rmtAddr, rmtCtrlAddr string, timeout time.Duration, gen Generator, csvW *csv.Writer) {
	kill := make(chan struct{})

	log.Printf("Resolving TCP addr for %v", rmtAddr)
	addr, err := net.ResolveTCPAddr("tcp", rmtAddr)
	if err != nil {
		log.Fatalf("Resolving IP addr, %v", err)
	}

	log.Printf("Creating TCP client with timeout %v", timeout)
	c := NewClient("tcp", addr, rmtCtrlAddr, timeout, gen)

	log.Printf("Starting to send")
	measures, err := c.Send(kill)
	if err != nil {
		log.Fatalf("Starting measurement client, %v", err)
	}

	log.Printf("Collecting measurements")
	csvW.Write([]string{
		"Time",
		"dT",
		"MessageSize",
		"Error",
	})
	for measure := range measures {
		var errStr string
		if measure.err != nil {
			errStr = measure.err.Error()
		}
		csvW.Write([]string{
			measure.t.String(),
			strconv.FormatInt(measure.dT.Nanoseconds(), 10),
			strconv.Itoa(measure.size),
			errStr,
		})
	}
	log.Printf("Done ranging over measurements")
}

func MeasureUDPClient(rmtAddr, rmtCtrlAddr string, timeout time.Duration, gen Generator, csvW *csv.Writer) {
	kill := make(chan struct{})

	log.Printf("Resolving UDP addr for %v", rmtAddr)
	addr, err := net.ResolveTCPAddr("tcp", rmtAddr)
	if err != nil {
		log.Fatalf("Resolving IP addr, %v", err)
	}

	log.Printf("Creating UDP client with timeout %v", timeout)
	c := NewClient("udp", addr, rmtCtrlAddr, timeout, gen)

	log.Printf("Starting to send")
	measures, err := c.Send(kill)
	if err != nil {
		log.Fatalf("Starting measurement client, %v", err)
	}

	log.Printf("Collecting measurements")
	csvW.Write([]string{
		"Time",
		"dT",
		"MessageSize",
		"Error",
	})
	for measure := range measures {
		var errStr string
		if measure.err != nil {
			errStr = measure.err.Error()
		}
		csvW.Write([]string{
			measure.t.String(),
			strconv.FormatInt(measure.dT.Nanoseconds(), 10),
			strconv.Itoa(measure.size),
			errStr,
		})
	}
	log.Printf("Done ranging over measurements")
}
