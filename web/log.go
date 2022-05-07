package web

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/oitel/tubelas/db"
)

const maxMessageCount = 10

func logHandler(w http.ResponseWriter, r *http.Request) {
	s := db.GlobalInstance()
	msgs, err := s.Load(maxMessageCount)
	if err != nil {
		log.Println("db.Storage.Load: ", err)
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
		return
	}
	render.JSON(w, r, msgs)
}
