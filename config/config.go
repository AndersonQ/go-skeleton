package config

import (
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"

	"github.com/AndersonQ/go-skeleton/constants"
)

// Config - environment variables are parsed to this struct
type Config struct {
	AppName    string `env:"APP_NAME" envDefault:"boilerplate"`
	Env        string `env:"ENV" envDefault:"env not set"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"debug"`
	LogOutput  string `env:"LOG_OUTPUT" envDefault:"console"`
	ServerPort int    `env:"PORT" envDefault:"8000"`

	// HTTP configurations
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" envDefault:"5s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"1s"`
	// RequestTimeout the timeout for the incoming request set on the request handler
	RequestTimeout time.Duration `env:"REQUEST_TIMEOUT" envDefault:"2s"`
	// WriteTimeout maximum time the server will handle a request before timing out writes of the response.
	// It must be bigger than RequestTimeout
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"4s"`

	// ShutdownTimeout the time the sever will wait server.Shutdown to return
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"6s"`

	// Opentracing
	OpentracingCollectorURL string `env:"OPENTRACING_COLLECTOR_URL" envDefault:"http://localhost:14268/api/traces"`
}

// Parse environment variables, returns (guess what?) and an error if an error occurs
func Parse() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}

// Logger returns a initialised zerolog.Logger
func (c Config) Logger() zerolog.Logger {
	logLevelOk := true
	logLevel, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
		logLevelOk = false
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimestampFieldName = constants.LogKeyTimestamp

	host, _ := os.Hostname()
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str(constants.LogKeyApp, c.AppName).
		Str(constants.LogKeyHost, host).
		Str(constants.LogKeyEnv, c.Env).
		Logger()

	if strings.ToUpper(c.LogOutput) == "CONSOLE" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if !logLevelOk {
		logger.Warn().Err(err).Msgf("%s is not a valid zerolog log level, defaulting to info", c.LogLevel)
	}

	return logger
}
