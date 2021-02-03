package server

import "net"

type Server interface {
	Connect() error
	StartListening()
	acceptClientConnection(conn net.Conn) *connectedClient
	serveClient(client *connectedClient)
}
