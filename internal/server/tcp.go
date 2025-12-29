package server

import (
	"net"
	"sync"
)

type TCPTransport struct {
	Network *net.TCPListener
	Addr    string
	users   map[string]*TCPPeer

	mut sync.Mutex
}

type TCPPeer struct {
	ip   string
	conn net.Conn
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		users: make(map[string]*TCPPeer),
		Addr:  addr,
	}
}

func (tp *TCPTransport) ListenAndAccept() error {
	server, err := net.Listen("tcp", tp.Addr)
	if err != nil {
		return err
	}

	tp.Network = server.(*net.TCPListener)

	for {
		conn, err := tp.Network.Accept()
		if err != nil {
			return err
		}

		peer := &TCPPeer{
			ip:   conn.RemoteAddr().String(),
			conn: conn,
		}
		tp.mut.Lock()
		tp.users[peer.ip] = peer
		tp.mut.Unlock()

		go func() {
			_ = handlePeerConnection(peer)
			_ = conn.Close()

			tp.mut.Lock()
			delete(tp.users, peer.ip)
			tp.mut.Unlock()
		}()
	}
}

func handlePeerConnection(peer *TCPPeer) error {
	// TODO: Implement read loop, encode and decode
	peer.Send([]byte("hi\n"))
	return nil
}

func (tp *TCPTransport) Close() error {
	if tp.Network == nil {
		return nil
	}

	return tp.Network.Close()
}
