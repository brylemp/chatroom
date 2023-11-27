package chatroom

import (
	"crypto/tls"
	"net"
)

func newListener(network, address string, tlsCfg *tls.Config) (net.Listener, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	if tlsCfg == nil {
		return listener, nil
	}

	return tls.NewListener(listener, tlsCfg), nil
}
