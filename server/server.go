package server

import "net"

type Server interface {
	Start() error
	Listen()
	Close()
	accept(conn net.Conn) *connectedClient
	serve(client *connectedClient)
	disconnect(client *connectedClient)
	broadcast(message interface{})
}
