package message

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strings"
)

var unknownMessageError = errors.New("unknown message")

type Reader struct {
	reader *bufio.Reader
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(reader),
	}
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
			ClientName: trimSuffix(clientName, ":"),
			Message:    trimSuffix(message, "\n"),
		}, nil

	default:
		log.Printf("Unknow message: %s", messageType)
	}

	return nil, unknownMessageError
}

func (r *Reader) ReadAll() ([]interface{}, error) {
	var commands []interface{}

	for {
		command, err := r.Read()

		if command != nil {
			commands = append(commands, command)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return commands, err
		}
	}

	return commands, nil
}