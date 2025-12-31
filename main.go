package main

import (
	"log"

	"github.com/ukpabik/mfs/internal/transport"
)

func main() {
	config := transport.NewTCPTransportConfig(":3000", transport.SimpleHandshake)
	tp := transport.NewTCPTransport(config)

	go func() {
		for rpc := range tp.Consume() {
			log.Printf("rpc from=%v payload=%s", rpc.From, string(rpc.Payload))
		}
		log.Printf("rpc channel closed")
	}()

	if err := tp.ListenAndAccept(); err != nil {
		log.Fatalf("server crashed")
	}
}
