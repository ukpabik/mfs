package server

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ukpabik/mfs/internal/files"
	"github.com/ukpabik/mfs/internal/transport"
)

func TestFileServer(t *testing.T) {
	tpCfg := transport.NewTCPTransportConfig(":3000", transport.SimpleHandshake)
	tp := transport.NewTCPTransport(tpCfg)

	fhCfg := files.NewFileHandlerConfig("./data-test")
	fh := files.NewFileHandler(fhCfg)

	cfg := NewFileServerConfig(tp, fh)

	require.NotEmpty(t, cfg.ID)
	require.Same(t, tp, cfg.Transport)
	require.Same(t, fh, cfg.FileHandler)
}
