package server

type HandshakeFunc func(Peer) error

func SimpleHandshake(p Peer) error {
	return nil
}
