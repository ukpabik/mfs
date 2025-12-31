package main

import (
	"fmt"

	"github.com/ukpabik/mfs/internal/server"
	"github.com/ukpabik/mfs/internal/transport"
)

var (
	addrOne = ":3000"
)

func main() {
	config := transport.NewTCPTransportConfig(addrOne, transport.SimpleHandshake)
	tp := transport.NewTCPTransport(config)

	rootDir := fmt.Sprintf("%s_network", addrOne)
	fsConfig := server.NewFileServerConfig(tp, rootDir)
	fs := server.NewFileServer(fsConfig)
	defer fs.Close()

	if err := fs.Start(); err != nil {
		panic(err)
	}
}
