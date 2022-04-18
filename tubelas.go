package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oitel/tubelas/web"
)

func main() {
	r := chi.NewRouter()
	r.Route("/", web.Route)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}
