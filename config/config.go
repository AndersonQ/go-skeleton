package config

import (
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"

	"github.com/AndersonQ/go-skeleton/constants"
)

// Config - environment variables are parsed to this struct
type Config struct {
	AppName    string `env:"APP_NAME" envDefault:"boilerplate"`
	ServerPort int    `env:"PORT" envDefault:"8000"`
	Env        string `env:"ENV" envDefault:"Development"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"debug"`
	LogOutput  string `env:"LOG_OUTPUT" envDefault:"console"`

	// NewRelicKey let it blank to disable newrelic
	NewRelicKey string `env:"NEWRELIC_KEY" envDefault:""`

	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" envDefault:"5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"3s"`
	// RequestTimeout the timeout for the incoming request
	RequestTimeout    time.Duration `env:"REQUEST_TIMEOUT" envDefault:"5s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"2s"`

	// ShutdownTimeout the time the sever will wait server.Shutdown to return
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"6s"`
}

// Parse environment variables, returns (guess what?) an error if an error occurs
func Parse() (Config, error) {
	confs := Config{}
	err := env.Parse(&confs)
	return confs, err
}

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

	if !logLevelOk {
		logger.Warn().Err(err).Msgf("%s is not a valid zerolog log level, defaulting to info", c.LogLevel)
	}

	return logger
}
