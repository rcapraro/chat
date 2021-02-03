package client

type Client interface {
	Connect(name string) error
	StartListening()
	SendMessage(message interface{}) error
	Close()
}