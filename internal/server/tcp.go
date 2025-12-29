package server

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransportConfig struct {
	Addr          string
	handshakeFunc HandshakeFunc
	// TODO: Add decoder here as well once implemented
}

type TCPTransport struct {
	config  TCPTransportConfig
	Network *net.TCPListener
	users   map[string]*TCPPeer

	mut sync.Mutex
}

type TCPPeer struct {
	conn net.Conn
}

func NewTCPTransportConfig(addr string, handshakeFunc HandshakeFunc) TCPTransportConfig {
	return TCPTransportConfig{
		Addr:          addr,
		handshakeFunc: handshakeFunc,
	}
}

func NewTCPTransport(config TCPTransportConfig) *TCPTransport {
	return &TCPTransport{
		users:  make(map[string]*TCPPeer),
		config: config,
	}
}

func (tp *TCPTransport) ListenAndAccept() error {
	server, err := net.Listen("tcp", tp.config.Addr)
	if err != nil {
		return err
	}

	tp.Network = server.(*net.TCPListener)

	for {
		conn, err := tp.Network.Accept()
		if err != nil {
			fmt.Printf("peer connection error: %v", err)
		}

		peer := &TCPPeer{
			conn: conn,
		}

		if tp.config.handshakeFunc != nil {
			if tp.config.handshakeFunc(peer); err != nil {
				_ = peer.Close()
				continue
			}
		}
		tp.mut.Lock()
		tp.users[peer.RemoteAddr().String()] = peer
		tp.mut.Unlock()

		go func() {
			_ = handlePeerConnection(peer)

			tp.mut.Lock()
			delete(tp.users, peer.RemoteAddr().String())
			tp.mut.Unlock()
		}()
	}
}

func handlePeerConnection(peer Peer) error {
	defer peer.Close()

	dec := NewGOBDecoder()
	enc := NewGOBEncoder()

	for {
		frame, err := readFrame(readerFromPeer{peer})
		if err != nil {
			return err
		}

		msg, err := dec.Decode(frame)
		if err != nil {
			return err
		}

		resp := Message{
			From:    peer.RemoteAddr(),
			Payload: append([]byte("ack: "), msg.Payload...),
		}

		out, err := enc.Encode(resp)
		if err != nil {
			return err
		}

		if err := writeFrame(writerToPeer{peer}, out); err != nil {
			return err
		}
	}
}

func (tp *TCPTransport) Close() error {
	if tp.Network == nil {
		return nil
	}

	return tp.Network.Close()
}

type readerFromPeer struct{ p Peer }

func (r readerFromPeer) Read(b []byte) (int, error) { return r.p.Read(b) }

type writerToPeer struct{ p Peer }

func (w writerToPeer) Write(b []byte) (int, error) {
	if err := w.p.Send(b); err != nil {
		return 0, err
	}
	return len(b), nil
}
