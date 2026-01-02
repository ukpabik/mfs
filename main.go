package main

import (
	"github.com/ukpabik/mfs/internal/manager"
	"github.com/ukpabik/mfs/internal/storage"
	"github.com/ukpabik/mfs/internal/transport"
)

func main() {
	node1 := storage.NewStorageNode(
		storage.NewStorageNodeConfig(
			transport.NewTCPTransport(transport.NewTCPTransportConfig(":3001", transport.SimpleHandshake)),
			"./data1",
		),
	)

	node2 := storage.NewStorageNode(
		storage.NewStorageNodeConfig(
			transport.NewTCPTransport(transport.NewTCPTransportConfig(":3002", transport.SimpleHandshake)),
			"./data2",
		),
	)

	node3 := storage.NewStorageNode(
		storage.NewStorageNodeConfig(
			transport.NewTCPTransport(transport.NewTCPTransportConfig(":3003", transport.SimpleHandshake)),
			"./data3",
		),
	)

	go func() { _ = node1.Start() }()
	go func() { _ = node2.Start() }()
	go func() { _ = node3.Start() }()

	mmTransport := transport.NewTCPTransport(
		transport.NewTCPTransportConfig(":3000", transport.SimpleHandshake),
	)
	mm := manager.NewMetadataManager("mm-1", mmTransport, []*storage.StorageNode{node1, node2, node3})
	defer mm.Close()
	defer node1.Close()
	defer node2.Close()
	defer node3.Close()

	if err := mm.Start(); err != nil {
		panic(err)
	}
}
