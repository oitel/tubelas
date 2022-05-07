package hub

type client struct {
	hub      *impl
	messages chan Message
}

func (cl client) Incoming() chan Message {
	return cl.messages
}

func (cl client) Publish(msg Message) {
	cl.hub.messages <- msg
}
