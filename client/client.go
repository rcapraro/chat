package client

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/rcapraro/chat/internal/message"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

type Client struct {
	Name          string
	MessagesChan  chan message.ServerMessage //Channel for the server messages
	conn          net.Conn                   //Client TCP socket connection
	messageWriter *message.Writer            //Writer to write messages on the conn
	messageReader *message.Reader            //Reader to read messages from the server
}

func NewClient() *Client {
	//Use the Seed function to initialize the default Source for Int in a non deterministic sequence of values.
	rand.Seed(time.Now().UnixNano())
	return &Client{
		Name:         fmt.Sprintf("%s %d", randomdata.FullName(randomdata.RandomGender), rand.Intn(99)),
		MessagesChan: make(chan message.ServerMessage),
	}
}

func (c *Client) Connect(name string) error {
	conn, err := net.Dial("tcp", "127.0.0.1:6697") //because it reminds me of IRC

	if err != nil {
		return err
	}
	log.Printf("Client connected to server / port 6697")

	c.conn = conn
	c.messageWriter = message.NewWriter(conn)
	c.messageReader = message.NewReader(conn)

	// Handle reconnection
	if name != "" {
		fmt.Print("\u001B[32m>\u001B[0m ")
		err = c.SendMessage(message.NameMessage{
			Name: c.Name,
		})
		if err != nil {
			log.Fatalf("Impossible to Send message to the server...exiting")
		}
	}

	return nil
}

func (c *Client) StartListening() {
	for {
		msg, err := c.messageReader.Read()

		if err != nil {
			log.Printf("Error while reading message: %v", err)
		}

		if ne, ok := err.(net.Error); ok && ne.Timeout() && ne.Temporary() || err == io.EOF {
			log.Printf("Network error: %v...trying to reconnect", err)
			time.Sleep(3 * time.Second)
			_ = c.Connect(c.Name)
			continue
		} else if err != nil {
			log.Fatalf("Unrecoverable error: %v...exiting", err)
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
