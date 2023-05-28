package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucian0ramos/image-golang/src/handlers"
	"github.com/lucian0ramos/image-golang/src/middlewares"
	"github.com/lucian0ramos/image-golang/src/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

func main() {
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro

	// logs (he hardcodeado la versión, lo suyo es automatizar con commit_id en git)
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("version", "1.0").
		Logger()

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.IDKey{}, uuid.New())

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	//router
	r := mux.NewRouter()
	r.HandleFunc("/manage-errors", middlewares.Chain(handlers.ManageErrors, middlewares.ValidateContentType(), middlewares.ValidateAuthorization())).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    *addr,
	}
	go serveHTTP(ctx, log, srv)

	<-quit

	log.Info().Interface("event_id", ctx.Value(models.IDKey{})).Msg("shutdown image-golang service")

	// Gracefully shutdown connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}

func serveHTTP(ctx context.Context, log zerolog.Logger, srv *http.Server) {
	log.Info().Interface("event_id", ctx.Value(models.IDKey{})).Msgf("image-golang started at %s", srv.Addr)
	err := srv.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Error().Caller().Interface("event_id", ctx.Value(models.IDKey{})).Err(err).Msg("starting Server listener failed")
	}
}
