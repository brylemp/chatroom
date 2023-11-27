package chatroom

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/brylemp/chatroom/pkg/util"
)

type userConnectionManager interface {
	AddUser(name string, conn net.Conn) error
	RemoveUser(name string)
	GetUser(name string) net.Conn
	GetUsers() map[string]net.Conn
}

type Chatroom struct {
	name    string
	network string
	address string

	tlsConfig *tls.Config
	logger    *slog.Logger

	userConnectionManager userConnectionManager

	errc     chan error
	messagec chan message
}

func New(ucm userConnectionManager, cro ...ChatroomOption) *Chatroom {
	cr := &Chatroom{
		name:      "Chatroom",
		network:   "tcp",
		address:   ":8080",
		tlsConfig: nil,
		logger:    slog.Default(),

		userConnectionManager: ucm,

		errc:     make(chan error),
		messagec: make(chan message),
	}

	for _, option := range cro {
		option(cr)
	}

	return cr
}

func (c *Chatroom) Start(ctx context.Context) error {
	listener, err := newListener(c.network, c.address, c.tlsConfig)
	if err != nil {
		return fmt.Errorf("Error listening: %w", err)
	}
	defer listener.Close()

	c.logger.Info(c.name+" started", "network", c.network, "address", c.address)

	// Handle messages
	go func() {
		for {
			msg := <-c.messagec

			err := c.broadcastMessage(msg)
			if err != nil {
				c.errc <- fmt.Errorf("Error broadcasting message: %w", err)
			}

			c.logger.Info("message sent", "type", msg.textType, "sender", msg.sender, "text", msg.text)
		}
	}()

	// Handle errors
	go func() {
		for {
			err := <-c.errc
			c.logger.Error(err.Error())
		}
	}()

	// Handle connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			c.errc <- fmt.Errorf("Error accepting connection: %w", err)
			continue
		}

		c.logger.Info("connection accepted", "addr", conn.RemoteAddr())

		err = util.SendMessage(conn, c.name)
		if err != nil {
			c.errc <- fmt.Errorf("Error sending message to user")
			continue
		}
		go c.handleConnection(conn)
	}
}

func (c *Chatroom) handleConnection(conn net.Conn) {
	defer conn.Close()

	userName, err := c.getUsername(conn)
	if err != nil {
		c.errc <- fmt.Errorf("Error getting username: %w", err)
		return
	}

	leaveFn, err := c.userJoins(userName, conn)
	if err != nil {
		c.errc <- fmt.Errorf("Error adding user: %w", err)
		return
	}
	defer leaveFn()

	for {
		msg, err := util.GetMessage(conn)
		if err == io.EOF {
			break
		}
		if err != nil {
			c.errc <- fmt.Errorf("Error getting user input: %w", err)
			break
		}

		c.messagec <- newMessage(userName, msg, chat)
	}
}

func (c *Chatroom) getUsername(conn net.Conn) (string, error) {
	var userName string

	for {
		input, err := util.GetMessage(conn)
		if err != nil {
			return "", fmt.Errorf("Error getting user input: %w", err)
		}
		userName = input

		userConn := c.userConnectionManager.GetUser(userName)
		if userConn == nil {
			break
		}

		err = util.SendMessage(conn, "Username already taken. Please try again.\n")
		if err != nil {
			return "", fmt.Errorf("Error sending message to user: %w", err)
		}
	}

	return userName, nil
}

func (c *Chatroom) broadcastMessage(msg message) error {
	toBroadcast := msg.getMessage()
	users := c.userConnectionManager.GetUsers()

	for _, conn := range users {
		if err := util.SendMessage(conn, toBroadcast); err != nil {
			return fmt.Errorf("Error sending message to user: %w", err)
		}
	}

	return nil
}

func (c *Chatroom) userJoins(name string, conn net.Conn) (func(), error) {
	joinMessage := name + " has joined"
	err := c.userConnectionManager.AddUser(name, conn)
	if err != nil {
		return nil, err
	}
	c.messagec <- newSystemMessage(joinMessage)

	return func() {
		leaveMessage := name + " has left"
		c.userConnectionManager.RemoveUser(name)
		c.messagec <- newSystemMessage(leaveMessage)
	}, nil
}
