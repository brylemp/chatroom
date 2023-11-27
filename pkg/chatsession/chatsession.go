package chatsession

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"strings"

	"github.com/brylemp/chatroom/pkg/util"
)

type GraphicUserInterfacer interface {
	Run(string) error
	SetInputReader(io.ReadWriter)
	GetOutputWriter() io.Writer
}

type ChatSession struct {
	network   string
	address   string
	tlsConfig *tls.Config

	gui GraphicUserInterfacer
}

func New(address string, gui GraphicUserInterfacer, options ...ChatSessionOption) *ChatSession {
	cs := &ChatSession{
		network: "tcp",
		address: address,

		gui: gui,
	}

	for _, option := range options {
		option(cs)
	}

	return cs
}

func (cs *ChatSession) Start(username string) error {
	conn, err := newDialer(cs.network, cs.address, cs.tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	chatroomName, err := util.GetMessage(conn)
	if err != nil {
		return err
	}

	if err := sendUsername(conn, username); err != nil {
		return err
	}

	cs.gui.SetInputReader(conn)
	go handleMessages(conn, cs.gui.GetOutputWriter())

	return cs.gui.Run(chatroomName)
}

func sendUsername(conn net.Conn, username string) error {
	if err := util.SendMessage(conn, username); err != nil {
		return err
	}

	isTaken, err := util.GetMessage(conn)
	if err != nil {
		return err
	}

	nameTaken := "Username already taken"
	if strings.Contains(isTaken, nameTaken) {
		return errors.New(nameTaken)
	}

	return nil
}

func handleMessages(conn net.Conn, w io.Writer) {
	for {
		message, err := util.GetMessage(conn)
		if err != nil {
			return
		}

		err = util.SendMessage(w, message)
		if err != nil {
			return
		}
	}
}
