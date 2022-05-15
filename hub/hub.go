package hub

import (
	"context"
	"log"
	"time"

	"github.com/oitel/tubelas/db"
	"github.com/oitel/tubelas/message"
)

type impl struct {
	clients    []client
	messages   chan message.Message
	register   chan client
	unregister chan Client
	killswitch chan struct{}
	storage    db.Storage
}

func newHub() Hub {
	return &impl{
		clients:    []client{},
		messages:   make(chan message.Message),
		register:   make(chan client),
		unregister: make(chan Client),
		killswitch: make(chan struct{}),
		storage:    db.GlobalInstance(), // FIXME: proper dependency injection
	}
}

const (
	storeTimeout = 30 * time.Second
)

func (h *impl) Run() error {
loop:
	for {
		select {
		case msg := <-h.messages:
			msg.Timestamp = time.Now().UTC().Unix()

			ctx, cancel := context.WithTimeout(context.Background(), storeTimeout)
			msg, err := h.storage.Store(ctx, msg)
			cancel()
			if err != nil {
				log.Println("db.Storage.Store: ", err)
				continue
			}

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
