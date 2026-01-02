package storage

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ukpabik/mfs/internal/files"
	"github.com/ukpabik/mfs/internal/transport"
)

// StorageNode is a replica node that stores files and executes operations.
// Each node listens on its own TCP port and processes requests from the MetadataManager.
type StorageNode struct {
	StorageNodeConfig

	stopCh chan struct{}
}

type StorageNodeConfig struct {
	ID string

	Transport   transport.Transport
	FileHandler *files.FileHandler
}

func NewStorageNode(config StorageNodeConfig) *StorageNode {
	return &StorageNode{
		StorageNodeConfig: config,
		stopCh:            make(chan struct{}),
	}
}

func NewStorageNodeConfig(transport transport.Transport, rootDir string) StorageNodeConfig {
	return StorageNodeConfig{
		ID:          generateServerID(),
		Transport:   transport,
		FileHandler: files.NewFileHandler(rootDir),
	}
}

// loop receives incoming RPCs, executes the requested operation, and sends back the result.
func (sn *StorageNode) loop() {
	for {
		select {
		case rpc := <-sn.Transport.Consume():
			message, err := ParseMessage(&rpc)
			if err != nil {
				log.Printf("error parsing message: %v", err)
				continue
			}
			log.Printf("from: %v, action: %v, filePath: %v, data: %v, size: %v",
				message.From, message.Action, message.FilePath, message.Data, message.Size)

			resp, err := sn.HandleMessage(message)
			if err != nil {
				_ = rpc.Peer.Send([]byte("ERROR: " + err.Error()))
				continue
			}

			_ = rpc.Peer.Send(resp)

		case <-sn.stopCh:
			return
		}
	}
}

func (sn *StorageNode) HandleMessage(msg StorageNodeMessage) ([]byte, error) {
	switch msg.Action {
	case Create:
		if err := sn.FileHandler.Create(msg.FilePath); err != nil {
			return nil, err
		}
		return []byte("OK"), nil

	case Delete:
		if err := sn.FileHandler.Delete(msg.FilePath); err != nil {
			return nil, err
		}
		return []byte("OK"), nil

	case Read:
		data, err := sn.FileHandler.Read(msg.FilePath, msg.Size)
		if err != nil {
			return nil, err
		}
		return data, nil
	case Write:
		n, err := sn.FileHandler.Write(msg.FilePath, bytes.NewReader(msg.Data))
		if err != nil {
			return nil, err
		}
		return []byte(fmt.Sprintf("OK: wrote %d bytes", n)), nil

	default:
		return nil, fmt.Errorf("unknown action: %d", msg.Action)
	}
}

func (sn *StorageNode) Start() error {
	go sn.loop()

	if err := sn.Transport.ListenAndAccept(); err != nil {
		return err
	}

	return nil
}

func (sn *StorageNode) Close() error {
	close(sn.stopCh)
	return sn.Transport.Close()
}
