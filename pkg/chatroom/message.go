package chatroom

type messageType = uint8

const (
	system messageType = iota
	chat
)

type message struct {
	sender   string
	text     string
	textType messageType
}

func newMessage(sender, text string, textType messageType) message {
	return message{
		sender:   sender,
		text:     text,
		textType: textType,
	}
}

func newSystemMessage(text string) message {
	return message{
		sender:   "system",
		text:     text,
		textType: system,
	}
}

func (m message) getMessage() string {
	if m.textType == system {
		return m.text + "\n"
	}

	return m.sender + ": " + m.text + "\n"
}
