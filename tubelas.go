package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oitel/tubelas/db"
	"github.com/oitel/tubelas/hub"
	"github.com/oitel/tubelas/web"
	"github.com/spf13/viper"
)

func main() {
	if err := loadConfig(); err != nil {
		log.Fatal("loadConfig: ", err)
	}

	addr := viper.GetString("listen")
	dbstring := viper.GetString("db")

	s := db.GlobalInstance()
	if err := s.Open(dbstring); err != nil {
		log.Fatal("db.Storage.Open: ", err)
	}
	defer s.Close()

	h := hub.GlobalInstance()
	go h.Run()

	r := chi.NewRouter()
	r.Route("/", web.Route)

	log.Println("Ready to serve.")
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}

	if err := h.Stop(); err != nil {
		log.Fatal("hub.Hub.Stop: ", err)
	}
}
