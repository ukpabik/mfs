package server

import "net"

type Transport interface {
	ListenAndAccept() error
}

type Peer interface {
	Send([]byte) error
	Read([]byte) (int, error)
	Close() error
	RemoteAddr() net.Addr
}
