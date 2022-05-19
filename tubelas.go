package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oitel/tubelas/db"
	"github.com/oitel/tubelas/hub"
	"github.com/oitel/tubelas/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

const (
	openTimeout = 120 * time.Second
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro

	if err := loadConfig(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to load config")
		os.Exit(1)
	}

	addr := viper.GetString("listen")
	dbstring := viper.GetString("db")

	s := db.GlobalInstance()
	{
		ctx, cancel := context.WithTimeout(context.Background(), openTimeout)
		defer cancel()

		if err := s.Open(ctx, dbstring); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to open DB")
			os.Exit(1)
		}
	}
	defer s.Close()

	h := hub.GlobalInstance()
	go h.Run()
	defer h.Stop()

	r := chi.NewRouter()
	r.Route("/", web.Route)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Info().
		Msg("Ready to serve")

	sigCtx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	gr, grCtx := errgroup.WithContext(sigCtx)
	gr.Go(func() error {
		return srv.ListenAndServe()
	})
	gr.Go(func() error {
		<-grCtx.Done()
		return srv.Shutdown(context.Background())
	})
	if err := gr.Wait(); err != http.ErrServerClosed {
		log.Error().
			Err(err).
			Msg("Failed to serve HTTP")
		os.Exit(1)
	}
}
