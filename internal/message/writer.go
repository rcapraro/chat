package message

import (
	"fmt"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) Write(message interface{}) error {
	switch m := message.(type) {
	case ClientMessage:
		return w.writeMessage(fmt.Sprintf("%s%s\n", SENT, m.Message))
	case NameMessage:
		return w.writeMessage(fmt.Sprintf("%s%s\n", NAME, m.Name))
	case ServerMessage:
		return w.writeMessage(fmt.Sprintf("%s%s:%s\n", RECEIVED, m.ClientName, m.Message))
	}
	return nil
}

func (w *Writer) writeMessage(message string) error {
	_, err := w.writer.Write([]byte(message))
	return err
}
