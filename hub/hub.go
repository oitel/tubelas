package hub

import (
	"context"
	"time"

	"github.com/oitel/tubelas/db"
	"github.com/oitel/tubelas/message"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

type impl struct {
	clients       []client
	messages      chan message.Message
	queue         chan message.Message
	register      chan client
	unregister    chan Client
	killswitch    chan struct{}
	storage       db.Storage
	connSemaphore *semaphore.Weighted
}

func newHub() Hub {
	storage := db.GlobalInstance() // FIXME: proper dependency injection
	var connSemaphore *semaphore.Weighted
	if storage.MaxConnCount() > 0 {
		connSemaphore = semaphore.NewWeighted(storage.MaxConnCount())
	}
	return &impl{
		clients:       []client{},
		messages:      make(chan message.Message),
		queue:         make(chan message.Message),
		register:      make(chan client),
		unregister:    make(chan Client),
		killswitch:    make(chan struct{}),
		storage:       storage,
		connSemaphore: connSemaphore,
	}
}

const (
	storeTimeout = 30 * time.Second
)

func (h *impl) Store(msg message.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), storeTimeout)
	defer cancel()

	if h.connSemaphore != nil {
		if err := h.connSemaphore.Acquire(ctx, 1); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to store message")
			return
		}
		defer h.connSemaphore.Release(1)
	}

	msg.Timestamp = time.Now().UTC().Unix()
	msg, err := h.storage.Store(ctx, msg)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to store message")
		return
	}

	h.queue <- msg
}

func (h *impl) Run() error {
loop:
	for {
		select {
		case msg := <-h.messages:
			go h.Store(msg)
		case msg := <-h.queue:
			for _, cl := range h.clients {
				cl.queue <- msg
			}
		case cl := <-h.register:
			h.clients = append(h.clients, cl)
			go cl.Listen()
		case cl := <-h.unregister:
			for i, hcl := range h.clients {
				if hcl == cl {
					h.clients = append(h.clients[:i], h.clients[i+1:]...)
					hcl.Close()
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
		hcl.Close()
	}
	close(h.messages)
	return nil
}

func (h *impl) Stop() error {
	close(h.killswitch)
	return nil
}

func (h *impl) Register() Client {
	cl := client{
		hub:      h,
		messages: make(chan message.Message, clientMessageQueueSize),
		queue:    make(chan message.Message),
	}
	h.register <- cl
	return cl
}

func (h *impl) Unregister(cl Client) {
	h.unregister <- cl
}
