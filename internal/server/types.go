package server

type Transport interface {
	ListenAndAccept() error
}

type Peer interface {
	Send([]byte) error
}
