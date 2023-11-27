package chatroom

import (
	"crypto/tls"
	"log/slog"
)

type ChatroomOption func(*Chatroom)

func WithName(name string) ChatroomOption {
	return func(cr *Chatroom) {
		cr.name = name
	}
}

func WithNetwork(network string) ChatroomOption {
	return func(cr *Chatroom) {
		cr.network = network
	}
}

func WithAddress(address string) ChatroomOption {
	return func(cr *Chatroom) {
		cr.address = address
	}
}

func WithLogger(logger *slog.Logger) ChatroomOption {
	return func(cr *Chatroom) {
		cr.logger = logger
	}
}

func WithTLS(cfg *tls.Config) ChatroomOption {
	return func(cr *Chatroom) {
		cr.tlsConfig = cfg
	}
}
