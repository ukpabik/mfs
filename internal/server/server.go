package server

import "github.com/ukpabik/mfs/internal/files"

type Server struct {
	transport *Transport
	handler   *files.FileHandler
}
