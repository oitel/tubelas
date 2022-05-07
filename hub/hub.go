package hub

import (
	"time"

	"github.com/oitel/tubelas/message"
)

type impl struct {
	clients    []client
	messages   chan message.Message
	register   chan client
	unregister chan Client
	killswitch chan struct{}
	counter    uint64
}

func newHub() Hub {
	return &impl{
		clients:    []client{},
		messages:   make(chan message.Message),
		register:   make(chan client),
		unregister: make(chan Client),
		killswitch: make(chan struct{}),
	}
}

func (h *impl) Run() error {
loop:
	for {
		select {
		case msg := <-h.messages:
			h.counter++
			msg.ID = h.counter
			msg.Timestamp = time.Now().UTC().Unix()

			for _, cl := range h.clients {
				cl.messages <- msg
			}
		case cl := <-h.register:
			h.clients = append(h.clients, cl)
		case cl := <-h.unregister:
			for i, hcl := range h.clients {
				if hcl == cl {
					h.clients = append(h.clients[:i], h.clients[i+1:]...)
					close(hcl.messages)
					break
				}
			}
		case _, ok := <-h.killswitch:
			if !ok {
				break loop
			}
		}
	}

	for _, hcl := range h.clients {
		close(hcl.messages)
	}
	return nil
}

func (h *impl) Stop() error {
	close(h.killswitch)
	return nil
}

func (h *impl) Register() Client {
	cl := client{
		hub:      h,
		messages: make(chan message.Message),
	}
	h.register <- cl
	return cl
}

func (h *impl) Unregister(cl Client) {
	h.unregister <- cl
}
