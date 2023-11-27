package util

import (
	"fmt"
	"io"
)

const messageBufferSize = 256

func GetMessage(rd io.Reader) (string, error) {
	b := make([]byte, messageBufferSize)
	n, err := rd.Read(b)
	if err != nil {
		return "", err
	}

	return string(b[:n]), nil
}

func SendMessage(w io.Writer, message string) error {
	_, err := fmt.Fprint(w, message)
	if err != nil {
		return err
	}

	return nil
}
