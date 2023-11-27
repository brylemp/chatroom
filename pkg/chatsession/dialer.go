package chatsession

import (
	"crypto/tls"
	"net"
)

func newDialer(network, address string, tlsCfg *tls.Config) (net.Conn, error) {
	if tlsCfg == nil {
		conn, err := net.Dial(network, address)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	conn, err := tls.Dial(network, address, tlsCfg)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
