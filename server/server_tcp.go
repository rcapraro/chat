package server

import (
	"github.com/rcapraro/chat/internal/message"
	"io"
	"log"
	"net"
	"sync"
)

type TcpServer struct {
	listener net.Listener       //Server is listening to TCP sockets
	clients  []*connectedClient //Slice of connected clients
	mutex    *sync.Mutex        //Synchronization primitive to update the list of connected clients in a goroutine context
}

type connectedClient struct {
	name   string
	conn   net.Conn        //Client TCP socket connection
	writer *message.Writer //Writer to write messages on the conn
}

func NewServer() *TcpServer {
	return &TcpServer{mutex: &sync.Mutex{}}
}

func (s *TcpServer) Start() error {
	l, err := net.Listen("tcp", ":6697") //because it reminds me of IRC
	if err != nil {
		return err
	} else {
		s.listener = l
	}
	log.Printf("Chat Server listening on port 6697")

	return err
}

func (s *TcpServer) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			client := s.accept(conn)
			go s.server(client)
		}
	}
}

func (s *TcpServer) accept(conn net.Conn) *connectedClient {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &connectedClient{
		conn:   conn,
		writer: message.NewWriter(conn),
	}

	s.clients = append(s.clients, client)

	log.Printf("Server acccepting new connection from client %v", conn.RemoteAddr())

	return client
}

func (s *TcpServer) server(client *connectedClient) {
	messageReader := message.NewReader(client.conn)
	defer s.disconnect(client)

	for {
		msg, err := messageReader.Read()

		if err != nil {
			if err == io.EOF {
				break //will trigger disconnect
			} else {
				log.Printf("Error while reading message: %v", err)
			}
		}

		if msg != nil {
			switch msgType := msg.(type) {
			case message.NameMessage:
				client.name = msgType.Name
			case message.ClientMessage:
				log.Printf("Receiving message from client %v (%s), broadcasting...", client.conn.RemoteAddr(), client.name)
				go s.broadcast(message.ServerMessage{
					ClientName: client.name,
					Message:    msgType.Message,
				})
			}
		}
	}
}

func (s *TcpServer) disconnect(client *connectedClient) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientAddr := client.conn.RemoteAddr()
	defer func() {
		err := client.conn.Close()
		if err != nil {
			log.Printf("Error while trying to disconnect client %v (%s)", clientAddr, client.name)
		}
	}()

	for i, c := range s.clients {
		if c == client {
			s.clients = removeClient(s.clients, i)
		}
	}

	log.Printf("Server disconnecting client %v (%s)", clientAddr, client.name)
}

func (s *TcpServer) broadcast(message interface{}) {
	for _, c := range s.clients {
		err := c.writer.Write(message)
		if err != nil {
			log.Printf("Error while broadcasting message to client %v (%s)", c.conn.RemoteAddr(), c.name)
		}
	}
}

func removeClient(clients []*connectedClient, i int) []*connectedClient {
	clients[i] = clients[len(clients)-1] //swap client to delete with last client
	return clients[:len(clients)-1]      //remove the last client
}

func (s *TcpServer) Close() {
	err := s.listener.Close()
	if err != nil {
		log.Fatal("Error while closing the server...exiting")
	}
}
