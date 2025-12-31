package transport

import "net"

type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}

type Peer interface {
	Send([]byte) error
	Close() error
	RemoteAddr() net.Addr
}
