package server

import (
	"encoding/binary"
	"errors"
	"io"
)

// Frame: [4 byte payload size | payload....]

const maxFrameSize = 8 << 20

func readFrame(r io.Reader) ([]byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}

	n := binary.BigEndian.Uint32(lenBuf[:])
	if n == 0 {
		return nil, errors.New("empty frame")
	}
	if n > maxFrameSize {
		return nil, errors.New("frame too large")
	}

	payload := make([]byte, n)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func writeFrame(w io.Writer, payload []byte) error {
	if len(payload) == 0 {
		return errors.New("empty frame")
	}
	if len(payload) > maxFrameSize {
		return errors.New("frame too large")
	}

	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(payload)))

	if _, err := w.Write(lenBuf[:]); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}
