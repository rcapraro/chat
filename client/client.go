package client

type Client interface {
	Connect(name string) error
	Start()
	Send(message interface{}) error
	Close()
}
