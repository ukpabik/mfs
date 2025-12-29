package server

import "net"

func (peer *TCPPeer) Send(data []byte) error {
	_, err := peer.conn.Write(data)
	return err
}

func (peer *TCPPeer) Close() error {
	if peer == nil {
		return nil
	}

	return peer.conn.Close()
}

func (peer *TCPPeer) Read(data []byte) (int, error) {
	if peer == nil {
		return 0, nil
	}

	return peer.conn.Read(data)
}
func (peer *TCPPeer) RemoteAddr() net.Addr {
	return peer.conn.RemoteAddr()
}
