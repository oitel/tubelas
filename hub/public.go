package hub

import "github.com/oitel/tubelas/message"

type Client interface {
	Incoming() chan message.Message
	Publish(text string)
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
