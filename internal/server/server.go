package server

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/ukpabik/mfs/internal/files"
	"github.com/ukpabik/mfs/internal/transport"
)

type FileServer struct {
	FileServerConfig

	fsMut sync.Mutex
}

type FileServerConfig struct {
	ID string

	Transport   transport.Transport
	FileHandler *files.FileHandler
}

func NewFileServer(config FileServerConfig) *FileServer {
	return &FileServer{FileServerConfig: config}
}

func NewFileServerConfig(transport transport.Transport, handler *files.FileHandler) FileServerConfig {
	return FileServerConfig{
		ID:          generateServerID(),
		Transport:   transport,
		FileHandler: handler,
	}
}

func generateServerID() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
