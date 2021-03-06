package web

import (
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

func debugRoute(r chi.Router) {
	r.Route("/pprof", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/profile", pprof.Profile)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)
	})
}
