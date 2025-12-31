package server

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/ukpabik/mfs/internal/files"
	"github.com/ukpabik/mfs/internal/transport"
)

type FileServer struct {
	FileServerConfig

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
			message, err := parseMessage(&rpc)
			if err != nil {
				log.Printf("error parsing message: %v", err)
				continue
			}
			log.Printf("from: %v, action: %v, filePath: %v, data: %v, size: %v",
				message.From, message.Action, message.FilePath, message.Data, message.Size)

			err = fs.handleMessage(message)
			if err != nil {
				_ = rpc.Peer.Send([]byte("ERROR: " + err.Error()))
				continue
			}

			_ = rpc.Peer.Send([]byte("OK"))

		case <-fs.stopCh:
			return
		}
	}
}

func (fs *FileServer) handleMessage(msg FileServerMessage) error {
	switch msg.Action {
	case Create:
		return fs.FileHandler.Create(msg.FilePath)
	case Delete:
		return fs.FileHandler.Delete(msg.FilePath)
	default:
		return nil
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
