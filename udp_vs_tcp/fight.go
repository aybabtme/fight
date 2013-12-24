package udp_vs_tcp

import (
	"encoding/csv"
	"log"
	"net"
	"strconv"
)

func MeasureTCPServer(lstnAddr, ctrlAddr, rmtAddr string, gen Generator, csvW *csv.Writer) {

	kill := make(chan struct{})

	ctrl := NewControlServer(ctrlAddr, rmtAddr)
	go func() {
		err := ctrl.Start(kill)
		if err != nil {
			log.Fatalf("Starting control server, %v", err)
		}
	}()
	defer ctrl.Close()

	srv := NewTCPServer(lstnAddr, gen)

	measures, err := srv.Measure(gen.Size(), kill)
	if err != nil {
		log.Fatalf("Starting measurement server (closeErr=%v), %v", ctrl.Close(), err)
	}

	for measure := range measures {
		csvW.Write([]string{
			measure.t.String(),
			strconv.Itoa(measure.size),
			string(measure.msg),
			measure.err.Error(),
		})
	}
}
