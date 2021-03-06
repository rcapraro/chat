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

type TcpClient struct {
	Name          string
	MessagesChan  chan message.ServerMessage //Channel for the server messages
	conn          net.Conn                   //Client TCP socket connection
	messageWriter *message.Writer            //Writer to write messages on the conn
	messageReader *message.Reader            //Reader to read messages from the server
}

func NewClient() *TcpClient {
	//Use the Seed function to initialize the default Source for Int in a non deterministic sequence of values.
	rand.Seed(time.Now().UnixNano())
	return &TcpClient{
		Name:         randomdata.FullName(randomdata.RandomGender),
		MessagesChan: make(chan message.ServerMessage),
	}
}

func (c *TcpClient) Connect(name string) error {
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
		err = c.Send(message.NameMessage{
			Name: c.Name,
		})
		if err != nil {
			log.Fatalf("Impossible to Send message to the server...exiting")
		}
	}

	return nil
}

func (c *TcpClient) Start() {
	for {
		msg, err := c.messageReader.Read()

		if ne, ok := err.(net.Error); ok && ne.Timeout() && ne.Temporary() || err == io.EOF {
			log.Printf("Network error: %v...trying to reconnect in 3s", err)
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

func (c *TcpClient) Send(message interface{}) error {
	return c.messageWriter.Write(message)
}

func (c *TcpClient) Close() {
	c.conn.Close()
}
