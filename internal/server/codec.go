package server

import (
	"bytes"
	"encoding/gob"
	"net"
)

type Encoder interface {
	Encode(Message) ([]byte, error)
}

type Decoder interface {
	Decode([]byte) error
}

type Message struct {
	From    net.Addr
	Payload []byte
}

type GOBDecoder struct{}
type GOBEncoder struct{}

func NewGOBDecoder() *GOBDecoder {
	gob.Register(Message{})
	return &GOBDecoder{}
}

func NewGOBEncoder() *GOBEncoder {
	gob.Register(Message{})
	return &GOBEncoder{}
}

func (gb *GOBDecoder) Decode(data []byte) (Message, error) {
	var msg Message
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&msg); err != nil {
		return Message{}, err
	}
	return msg, nil
}

func (ge *GOBEncoder) Encode(msg Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
