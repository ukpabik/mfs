package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/ukpabik/mfs/internal/transport"
)

// Action represents a file operation type.
type Action byte

const (
	Create Action = 0x00
	Delete Action = 0x01
	Read   Action = 0x02
	Write  Action = 0x03
)

const MaxPayloadSize = 10 * 1024 * 1024 // Max size is 10MiB

// StorageNodeMessage represents a request to a storage node.
type StorageNodeMessage struct {
	From     net.Addr
	Action   Action
	FilePath string

	// Read/Write fields
	Size uint64
	Data []byte
}

// ParseMessage decodes an RPC payload into a StorageNodeMessage.
// Binary format: [action:1][pathLen:8][path:variable][optional params]
func ParseMessage(rpc *transport.RPC) (StorageNodeMessage, error) {
	if len(rpc.Payload) > MaxPayloadSize {
		return StorageNodeMessage{}, fmt.Errorf("payload too large: %d bytes", len(rpc.Payload))
	}

	r := bytes.NewReader(rpc.Payload)
	fsm := StorageNodeMessage{From: rpc.From}

	actionBuf := make([]byte, 1)
	if _, err := r.Read(actionBuf); err != nil {
		return StorageNodeMessage{}, fmt.Errorf("read action: %w", err)
	}
	fsm.Action = Action(actionBuf[0])

	var pathLen uint64
	if err := binary.Read(r, binary.BigEndian, &pathLen); err != nil {
		return StorageNodeMessage{}, fmt.Errorf("read filepath length: %w", err)
	}

	if pathLen == 0 || pathLen > 256 {
		return StorageNodeMessage{}, fmt.Errorf("invalid filepath length: %d", pathLen)
	}

	pathBuf := make([]byte, pathLen)
	if _, err := io.ReadFull(r, pathBuf); err != nil {
		return StorageNodeMessage{}, fmt.Errorf("read filepath: %w", err)
	}
	fsm.FilePath = string(pathBuf)

	switch fsm.Action {
	case Create, Delete:

	case Read:
		if err := binary.Read(r, binary.BigEndian, &fsm.Size); err != nil {
			return StorageNodeMessage{}, fmt.Errorf("read size param: %w", err)
		}

	case Write:
		var dataLen uint64
		if err := binary.Read(r, binary.BigEndian, &dataLen); err != nil {
			return StorageNodeMessage{}, fmt.Errorf("read data length: %w", err)
		}

		if dataLen == 0 || dataLen > MaxPayloadSize {
			return StorageNodeMessage{}, fmt.Errorf("invalid data length: %d", dataLen)
		}

		fsm.Data = make([]byte, dataLen)
		if _, err := io.ReadFull(r, fsm.Data); err != nil {
			return StorageNodeMessage{}, fmt.Errorf("read data: %w", err)
		}

	default:
		return StorageNodeMessage{}, fmt.Errorf("unknown action: %d", fsm.Action)
	}

	return fsm, nil
}
