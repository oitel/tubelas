package hub

type Message []byte

type Client interface {
	Incoming() chan Message
	Publish(msg Message)
}

type Hub interface {
	Run() error
	Stop() error

	Register() Client
	Unregister(cl Client)
}

func NewHub() Hub {
	return newHub()
}

var global Hub

func GlobalInstance() Hub {
	if global == nil {
		global = newHub()
	}
	return global
}
