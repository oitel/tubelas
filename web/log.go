package web

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/oitel/tubelas/db"
	"github.com/rs/zerolog/log"
)

const maxMessageCount = 10

func logHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := db.GlobalInstance()
	msgs, err := s.Load(ctx, maxMessageCount)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to load messages")
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
		return
	}
	render.JSON(w, r, msgs)
}
