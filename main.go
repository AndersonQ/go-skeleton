package main

import (
	"compress/flate"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/AndersonQ/go-skeleton/config"
	"github.com/AndersonQ/go-skeleton/handlers"
	"github.com/AndersonQ/go-skeleton/middlewares"
)

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)  // a.k.a ctrl+C
	signal.Notify(signalChan, syscall.SIGTERM) // a.k.a kill

	// when closed the program should exit
	idleConnsClosed := make(chan bool)

	cfg, err := config.Parse()
	if err != nil {
		panic("could not parse environment variables: " + err.Error())
	}

	logger := cfg.Logger()
	newrelicApp, err := initNewrelic(cfg, logger)
	if err != nil {
		logger.Warn().Err(err).Msg("could  init newrelic agent")
	}

	router := initRouter(cfg, newrelicApp, logger)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,

		IdleTimeout:       cfg.IdleTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	// handle graceful shutdown in another goroutine
	go gracefullShutdown(signalChan, idleConnsClosed, cfg, server, newrelicApp, logger)

	logger.Info().Msgf("staring server on :%d", cfg.ServerPort)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Err(err).Msg("HTTP server error")
	} else {
		logger.Info().Msg("shutting server down...")
	}

	<-idleConnsClosed
}

func initRouter(cfg config.Config, newrelicApp newrelic.Application, logger zerolog.Logger) *chi.Mux {
	router := chi.NewRouter()
	// TODO(Anderson): create a tracking ID middleware
	router.Use(hlog.NewHandler(logger))
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Compress(flate.BestSpeed))
	router.Use(middlewares.JsonResponse)
	router.Use(middlewares.RequestLogWrapper)
	router.Use(middlewares.TimeoutWrapper(cfg.RequestTimeout))

	router.Get(newrelic.WrapHandleFunc(newrelicApp, "/live", handlers.NewLivenessHandler()))
	router.Get(newrelic.WrapHandleFunc(newrelicApp, "/ready", handlers.NewReadinessHandler()))
	router.Get(newrelic.WrapHandleFunc(newrelicApp, "/slow", handlers.NewSlowHandler(
		cfg.RequestTimeout+time.Millisecond*5)))

	return router
}

func initNewrelic(cfg config.Config, logger zerolog.Logger) (newrelic.Application, error) {
	ncfg := newrelic.NewConfig(cfg.AppName, cfg.NewRelicKey)

	// Disable communication with newrelic, see https://github.com/newrelic/go-agent/blob/master/config.go#L27
	if len(cfg.NewRelicKey) == 0 {
		ncfg.Enabled = false
		logger.Warn().Msg("Disabling NewRelic as license key is empty")
	}

	app, err := newrelic.NewApplication(ncfg)

	return app, err
}

func gracefullShutdown(
	signalChan chan os.Signal,
	idleConnsClosed chan bool,
	cfg config.Config,
	server *http.Server,
	newrelicApp newrelic.Application,
	logger zerolog.Logger) {

	sig := <-signalChan
	logger.Info().Msgf("received signal: %q, starting graceful shutdown...", sig.String())

	ctx, done := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer done() // avoid a context leak

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("error during gracefully shutdown")
	}

	deadline, _ := ctx.Deadline()

	if newrelicApp != nil {
		logger.Info().Msg("shunting down Newrelic")
		newrelicApp.Shutdown(deadline.Sub(time.Now()))
	} else {
		logger.Info().Msg("newrelic not initialised, nothing to shutdown")
	}

	logger.Info().Msg("Gracefully shutdown finished")

	close(idleConnsClosed)
}
