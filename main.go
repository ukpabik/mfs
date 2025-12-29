package main

import (
	"log"

	"github.com/ukpabik/mfs/internal/server"
)

func main() {
	config := server.NewTCPTransportConfig(":3000", server.SimpleHandshake)
	tp := server.NewTCPTransport(config)

	if err := tp.ListenAndAccept(); err != nil {
		log.Fatalf("server crashed")
	}
}
