package hub

import "github.com/oitel/tubelas/message"

type client struct {
	hub      *impl
	messages chan message.Message
}

func (cl client) Incoming() chan message.Message {
	return cl.messages
}

func (cl client) Publish(text string) {
	msg := message.Message{
		Text: text,
	}
	cl.hub.messages <- msg
}
