package main

import (
	"github.com/rcapraro/chat/server"
	"log"
)

func main() {
	s := server.NewServer()
	defer s.Close()
	err := s.Start()
	if err != nil {
		log.Fatalf("Impossible to Start the server...exiting")
	}
	s.Listen()
}
