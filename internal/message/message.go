package message

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	SENT     = "[SENT]"
	NAME     = "[NAME]"
	RECEIVED = "[RECEIVED]"
)

var unknownMessageError = errors.New("unknown message")

type ClientMessage struct {
	Message string
}

type NameMessage struct {
	Name string
}

type ServerMessage struct {
	ClientName string
	Message    string
}

type Writer struct {
	writer io.Writer
}

type Reader struct {
	reader *bufio.Reader
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(reader),
	}
}

func (w *Writer) Write(message interface{}) error {
	switch m := message.(type) {
	case ClientMessage:
		return w.writeMessage(fmt.Sprintf("%s %s\n", SENT, m.Message))
	case NameMessage:
		return w.writeMessage(fmt.Sprintf("%s %s\n", NAME, m.Name))
	case ServerMessage:
		return w.writeMessage(fmt.Sprintf("%s %s: %s\n", RECEIVED, m.ClientName, m.Message))
	}
	return nil
}

func (w *Writer) writeMessage(message string) error {
	_, err := w.writer.Write([]byte(message))
	return err
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func (r *Reader) Read() (interface{}, error) {
	messageType, err := r.reader.ReadString(']')

	if err != nil {
		return nil, err
	}

	switch messageType {
	case SENT:
		message, err := r.reader.ReadString('\n')

		if err != nil {
			return nil, err
		}

		return ClientMessage{trimSuffix(message, "\n")}, nil

	case NAME:
		name, err := r.reader.ReadString('\n')

		if err != nil {
			return nil, err
		}

		return NameMessage{trimSuffix(name, "\n")}, nil

	case RECEIVED:
		clientName, err := r.reader.ReadString(':')
		if err != nil {
			return nil, err
		}

		message, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return ServerMessage{
			ClientName: trimSuffix(clientName, " "),
			Message:    trimSuffix(message, "\n"),
		}, nil

	default:
		log.Printf("Unknow message: %s", messageType)
	}

	return nil, unknownMessageError
}
