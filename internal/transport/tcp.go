package transport

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransportConfig struct {
	Addr          string
	handshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportConfig
	Network *net.TCPListener
	users   map[string]*TCPPeer

	mut        sync.Mutex
	rpcChannel chan RPC
}

type TCPPeer struct {
	conn net.Conn
}

const chanSize = 1024

func NewTCPTransportConfig(addr string, handshakeFunc HandshakeFunc) TCPTransportConfig {
	return TCPTransportConfig{
		Addr:          addr,
		handshakeFunc: handshakeFunc,
		Decoder:       SimpleDecoder{},
	}
}

func NewTCPTransport(config TCPTransportConfig) *TCPTransport {
	return &TCPTransport{
		users:              make(map[string]*TCPPeer),
		TCPTransportConfig: config,
		rpcChannel:         make(chan RPC, chanSize),
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
			fmt.Printf("peer connection error: %v", err)
			continue
		}

		peer := &TCPPeer{
			conn: conn,
		}

		if tp.handshakeFunc != nil {
			if err := tp.handshakeFunc(peer); err != nil {
				_ = peer.Close()
				continue
			}
		}
		tp.mut.Lock()
		tp.users[peer.RemoteAddr().String()] = peer
		tp.mut.Unlock()

		go func() {
			_ = tp.handleConnection(peer)

			tp.mut.Lock()
			delete(tp.users, peer.RemoteAddr().String())
			tp.mut.Unlock()
		}()
	}
}

func (tp *TCPTransport) handleConnection(peer *TCPPeer) error {
	defer peer.Close()

	for {
		rpc := RPC{}

		err := tp.Decoder.Decode(peer.conn, &rpc)
		if err != nil {
			return err
		}
		rpc.From = peer.RemoteAddr()

		tp.rpcChannel <- rpc
	}
}

func (tp *TCPTransport) Consume() <-chan RPC {
	return tp.rpcChannel
}

func (tp *TCPTransport) Close() error {
	return tp.Network.Close()
}
