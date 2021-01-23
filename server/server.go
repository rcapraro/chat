package server

import (
	"github.com/rcapraro/chat/internal/message"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	listener net.Listener       //Server is listening to TCP sockets
	clients  []*connectedClient //Slice of connected clients
	mutex    *sync.Mutex        //Synchronization primitive to update the list of connected clients in a goroutine context
}

type connectedClient struct {
	name   string
	conn   net.Conn        //Client TCP socket connection
	writer *message.Writer //Writer to write messages on the conn
}

func NewServer() *Server {
	return &Server{mutex: &sync.Mutex{}}
}

func (s *Server) Connect() error {
	l, err := net.Listen("tcp", ":6697") //because it reminds me of IRC
	if err != nil {
		return err
	} else {
		s.listener = l
	}
	log.Printf("Chat Server listening on port 6697")

	return err
}

func (s *Server) StartListening() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			client := s.acceptClientConnection(conn)
			go s.serveClient(client)
		}
	}
}

func (s *Server) acceptClientConnection(conn net.Conn) *connectedClient {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &connectedClient{
		conn:   conn,
		writer: message.NewWriter(conn),
	}

	s.clients = append(s.clients, client)

	log.Printf("Server acccepting new connection from client %s", conn.RemoteAddr().String())

	return client
}

func (s *Server) serveClient(client *connectedClient) {
	messageReader := message.NewReader(client.conn)
	defer s.disconnectClient(client)

	for {
		msg, err := messageReader.Read()

		if err != nil && err != io.EOF {
			log.Printf("Error while reading message: %v", err)
		}

		if msg != nil {
			switch msgType := msg.(type) {
			case message.NameMessage:
				client.name = msgType.Name
			case message.ClientMessage:
				log.Printf("Receiving message from client %s, broadcasting...", client.name)
				go s.broadcastMessage(message.ServerMessage{
					ClientName: client.name,
					Message:    msgType.Message,
				})
			}
		}
	}
}

func (s *Server) disconnectClient(client *connectedClient) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientAddr := client.conn.RemoteAddr().String()
	defer func() {
		err := client.conn.Close()
		if err != nil {
			log.Printf("Error while trying to disconnect client %s", clientAddr)
		}
	}()

	for i, c := range s.clients {
		if c == client {
			s.clients = removeClient(s.clients, i)
		}
	}

	log.Printf("Server disconnecting client %s", clientAddr)
}

func (s *Server) broadcastMessage(message interface{}) {
	for _, client := range s.clients {
		err := client.writer.Write(message)
		if err != nil {
			log.Printf("Error while broadcasting message to client %v", client.conn.RemoteAddr().String())
		}
	}
}

func removeClient(clients []*connectedClient, i int) []*connectedClient {
	clients[i] = clients[len(clients)-1] //swap client to delete with last client
	return clients[:len(clients)-1]      //remove the last client
}

func (s *Server) Close() {
	err := s.listener.Close()
	if err != nil {
		log.Fatal("Error while closing the server...exiting")
	}
}
