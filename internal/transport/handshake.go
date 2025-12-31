package transport

type HandshakeFunc func(Peer) error

func SimpleHandshake(p Peer) error {
	return nil
}
