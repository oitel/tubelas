package web

import "github.com/go-chi/chi/v5"

func Route(r chi.Router) {
	r.Get("/", helloHandler)
}
