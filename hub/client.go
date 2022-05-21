package hub

import (
	"time"

	"github.com/oitel/tubelas/message"
)

const (
	clientMessageQueueSize = 128
)

type client struct {
	hub      *impl
	messages chan message.Message

	queue     chan message.Message
	discarded bool
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

func (cl client) Listen() {
	for msg := range cl.queue {
		if len(cl.messages) == cap(cl.messages) {
			cl.discarded = true
		} else {
			if cl.discarded {
				for range cl.messages {
					// flush message queue
				}
				cl.messages <- message.Message{
					Timestamp: time.Now().UTC().Unix(),
					Text:      "Message queue was discarded",
				}
				cl.discarded = false
			}
			cl.messages <- msg
		}
	}
	close(cl.messages)
}

func (cl client) Close() {
	close(cl.queue)
}
