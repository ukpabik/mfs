package main

import (
	"log"

	"github.com/ukpabik/mfs/internal/server"
)

func main() {
	config := server.NewTCPTransportConfig(":3000", server.SimpleHandshake)
	tp := server.NewTCPTransport(config)

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
