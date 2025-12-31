package transport

import (
	"io"
	"net"
)

var maxBufferSize = 1028

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type RPC struct {
	From    net.Addr
	Peer    Peer
	Payload []byte
}

type SimpleDecoder struct{}

func (def SimpleDecoder) Decode(r io.Reader, msg *RPC) error {
	buf := make([]byte, maxBufferSize)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}
