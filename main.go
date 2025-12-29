package main

import (
	"log"

	"github.com/ukpabik/mfs/internal/server"
)

func main() {
	tp := server.NewTCPTransport(":3000")

	if err := tp.ListenAndAccept(); err != nil {
		log.Fatalf("server crashed")
	}
}
