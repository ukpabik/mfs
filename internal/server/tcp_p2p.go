package server

func (peer *TCPPeer) Send(data []byte) error {
	// TODO: Ensure data is correctly formatted --> create format checking func
	_, err := peer.conn.Write(data)
	return err
}
