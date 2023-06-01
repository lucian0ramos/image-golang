package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justinas/alice"
	"github.com/lucian0ramos/image-golang/src/handlers"
	"github.com/lucian0ramos/image-golang/src/middlewares"
	"github.com/lucian0ramos/image-golang/src/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	inFlightGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "in_flight_requests",
			Help: "A gauge of requests currently being served by the wrapped handler.",
		},
	)

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "A counter for requests to the wrapped handler.",
			ConstLabels: map[string]string{
				"version": version,
			},
		},
		[]string{"code", "method"},
	)

	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
			ConstLabels: map[string]string{
				"version": version,
			},
		},
		[]string{"code", "method"},
	)

	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
			ConstLabels: map[string]string{
				"version": version,
			},
		},
		[]string{"code", "method"},
	)
	version string
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
	c := alice.New(hlog.NewHandler(log), hlog.AccessHandler(accessLogger))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	r.HandleFunc("/manage-errors", middlewares.Chain(handlers.ManageErrors, middlewares.ValidateContentType(), middlewares.ValidateAuthorization())).Methods("POST")

	srv := &http.Server{
		Handler: c.Then(promRequestHandler(r)),
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

func accessLogger(r *http.Request, status, size int, dur time.Duration) {
	hlog.FromRequest(r).Info().
		Str("host", r.Host).
		Int("status", status).
		Str("url", r.RequestURI).
		Str("method", r.Method).
		Int("size", size).
		Dur("duration_ms", dur).
		Msg("request")
}

func promRequestHandler(handler http.Handler) http.Handler {
	return promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration,
			promhttp.InstrumentHandlerCounter(counter,
				promhttp.InstrumentHandlerResponseSize(responseSize, handler),
			),
		),
	)
}
