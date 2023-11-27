package chatsession

import "crypto/tls"

type ChatSessionOption func(*ChatSession)

func WithNetwork(network string) ChatSessionOption {
	return func(cs *ChatSession) {
		cs.network = network
	}
}

func WithTLS(tlsCfg *tls.Config) ChatSessionOption {
	return func(cs *ChatSession) {
		cs.tlsConfig = tlsCfg
	}
}
