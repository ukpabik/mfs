package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ukpabik/mfs/internal/transport"
)

func TestFileServer(t *testing.T) {
	tpCfg := transport.NewTCPTransportConfig(":3000", transport.SimpleHandshake)
	tp := transport.NewTCPTransport(tpCfg)

	cfg := NewStorageNodeConfig(tp, "./data-test")

	require.NotEmpty(t, cfg.ID)
	require.Same(t, tp, cfg.Transport)
}
