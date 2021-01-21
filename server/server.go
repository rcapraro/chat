package server

import (
	"github.com/rcapraro/chat/internal/message"
	"log"
	"net"
	"sync"
)

type Server struct {
	listener net.Listener
	clients  []*client
	mutex *sync.Mutex
}

type client struct {
	conn   net.Conn
	name   string
	writer *message.Writer
}

func NewServer() *Server {
	return &Server{mutex: &sync.Mutex{}}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", "127.0.0.1:8787")
	if err != nil {
		return err
	} else {
		s.listener = l
	}
	log.Printf("Server listening on 127.0.0.1 / port 8787")

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			client:=s.acceptClientConnection(conn)
			go s.handleClientMessage(client)
		}
	}
}

func (s *Server) acceptClientConnection(conn net.Conn) *client {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn: conn,
		writer: message.NewWriter(conn),
	}

	s.clients = append(s.clients, client)

	log.Printf("Server acccepting new connection from client %v", conn.RemoteAddr().String())

	return client
}

func (s *Server) handleClientMessage(client *client) {
	//TODO
}



