package message

const (
	SENT     = "[SENT]"
	NAME     = "[NAME]"
	RECEIVED = "[RECEIVED]"
)

type ClientMessage struct {
	Message string
}

type NameMessage struct {
	Name string
}

type ServerMessage struct {
	ClientName string
	Message    string
}
