package transport

import "net"

func (peer *TCPPeer) Send(data []byte) error {
	_, err := peer.conn.Write(data)
	return err
}

func (peer *TCPPeer) Close() error {
	return peer.conn.Close()
}

func (peer *TCPPeer) Read(data []byte) (int, error) {
	return peer.conn.Read(data)
}
func (peer *TCPPeer) RemoteAddr() net.Addr {
	return peer.conn.RemoteAddr()
}
