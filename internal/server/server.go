package server

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"sync"

	"github.com/ukpabik/mfs/internal/files"
	"github.com/ukpabik/mfs/internal/transport"
)

type FileServer struct {
	FileServerConfig

	fsMut  sync.Mutex
	stopCh chan struct{}
}

type FileServerConfig struct {
	ID string

	Transport   transport.Transport
	FileHandler *files.FileHandler
}

func NewFileServer(config FileServerConfig) *FileServer {
	return &FileServer{
		FileServerConfig: config,
		stopCh:           make(chan struct{}),
	}
}

func NewFileServerConfig(transport transport.Transport, rootDir string) FileServerConfig {
	return FileServerConfig{
		ID:          generateServerID(),
		Transport:   transport,
		FileHandler: files.NewFileHandler(rootDir),
	}
}

func generateServerID() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}

func (fs *FileServer) loop() {
	for {
		select {
		case rpc := <-fs.Transport.Consume():
			// TODO: Convert the RPC to a specific type of message and handle it
			log.Printf("rpc from=%v payload=%s", rpc.From, string(rpc.Payload))
		case <-fs.stopCh:
			return
		}
	}
}

func (fs *FileServer) Start() error {
	go fs.loop()

	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	return nil
}

func (fs *FileServer) Close() error {
	close(fs.stopCh)
	return fs.Transport.Close()
}
