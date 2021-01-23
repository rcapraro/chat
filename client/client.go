package client

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/rcapraro/chat/internal/message"
	"io"
	"log"
	"net"
)

type Client struct {
	Name          string
	MessagesChan  chan message.ServerMessage //Channel for the server messages
	conn          net.Conn                   //Client TCP socket connection
	messageWriter *message.Writer            //Writer to write messages on the conn
	messageReader *message.Reader            //Reader to read messages from the server
}

func NewClient() *Client {
	return &Client{
		Name: randomdata.FullName(randomdata.RandomGender),
		MessagesChan: make(chan message.ServerMessage),
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", "127.0.0.1:6697") //because it reminds me of IRC

	if err != nil {
		return err
	}
	log.Printf("Client connected to server / port 6697")

	c.conn = conn
	c.messageWriter = message.NewWriter(conn)
	c.messageReader = message.NewReader(conn)

	return nil
}

func (c *Client) StartListening() {
	for {
		msg, err := c.messageReader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error while reading message: %v", err)
		}

		if msg != nil {
			switch msgType := msg.(type) {
			case message.ServerMessage:
				c.MessagesChan <- msgType
			default:
				log.Printf("Unknown message type")
			}
		}
	}
}

func (c *Client) SendMessage(message interface{}) error {
	return c.messageWriter.Write(message)
}

func (c *Client) Close() {
	c.conn.Close()
}


