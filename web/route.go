package web

import "github.com/go-chi/chi/v5"

func Route(r chi.Router) {
	r.Route("/debug", debugRoute)

	r.Get("/", helloHandler)
	r.Get("/log", logHandler)
	r.Get("/ws", wsHandler)
}
