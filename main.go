package main

import (
	"compress/flate"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	bllogger "github.com/blacklane/go-libs/logger"
	blmiddleware "github.com/blacklane/go-libs/logger/middleware"
	bltrackingmiddleware "github.com/blacklane/go-libs/tracking/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"

	"github.com/AndersonQ/go-skeleton/config"
	"github.com/AndersonQ/go-skeleton/handlers"
	"github.com/AndersonQ/go-skeleton/middlewares"
	"github.com/AndersonQ/go-skeleton/tracing"
)

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	// when closed the program should exit
	idleConnsClosed := make(chan bool)

	cfg, err := config.Parse()
	if err != nil {
		panic("could not parse environment variables: " + err.Error())
	}

	logger := cfg.Logger()

	tracer, tracerCloser := tracing.NewJaegerTracer(cfg, logger)

	router := initRouter(cfg, logger, tracer)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,

		IdleTimeout:       cfg.IdleTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	// handle graceful shutdown in another goroutine
	go gracefulShutdown(signalChan, idleConnsClosed, cfg, server, logger, tracerCloser)

	logger.Info().Msgf("staring server on :%d", cfg.ServerPort)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Err(err).Msg("HTTP server error")
	} else {
		logger.Info().Msg("shutting server down...")
	}

	<-idleConnsClosed
}

func initRouter(cfg config.Config, logger bllogger.Logger, tracer opentracing.Tracer) *chi.Mux {
	router := chi.NewRouter()
	router.Use(bltrackingmiddleware.TrackingID)
	router.Use(blmiddleware.HTTPAddLogger(logger))
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Compress(flate.BestSpeed))
	router.Use(middlewares.JsonResponse)
	router.Use(blmiddleware.HTTPRequestLogger([]string{}))
	router.Use(middlewares.TimeoutWrapper(cfg.RequestTimeout))

	router.Get(middlewares.AddOpentracing("/live", tracer, handlers.NewLivenessHandler()))
	router.Get(middlewares.AddOpentracing("/ready", tracer, handlers.NewReadinessHandler()))
	router.Get(middlewares.AddOpentracing("/slow", tracer, handlers.NewSlowHandler(
		cfg.RequestTimeout+time.Millisecond*5)))

	router.Get(middlewares.AddOpentracing("/trace", tracer, handlers.NewOpentracing()))

	return router
}

func gracefulShutdown(
	signalChan chan os.Signal,
	idleConnsClosed chan bool,
	cfg config.Config,
	server *http.Server,
	logger zerolog.Logger,
	tracer io.Closer) {

	sig := <-signalChan
	logger.Info().Msgf("received signal: %q, starting graceful shutdown...", sig.String())

	ctx, done := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer done() // avoid a context leak

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("error during server shutdown")
	}

	if err := tracer.Close(); err != nil {
		logger.Error().Err(err).Msg("error during tracer shutdown")
	}

	logger.Info().Msg("Gracefully shutdown finished")

	close(idleConnsClosed)
}
