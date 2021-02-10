package tracing

import (
	"errors"
	"fmt"
	"io"

	"github.com/blacklane/go-libs/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/AndersonQ/go-skeleton/config"
)

// NewJaegerTracer returns a Jeager implementation of opentracing.Tracer
func NewJaegerTracer(cfg config.Config, logger logger.Logger) (opentracing.Tracer, io.Closer) {
	jaegerCfg := jaegerconfig.Configuration{
		ServiceName: cfg.AppName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: cfg.OpentracingCollectorURL,
		},
	}

	// Example metrics factory. Use github.com/uber/jaeger-lib/metrics to bind
	// to real metric framework.
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := jaegerCfg.NewTracer(
		jaegerconfig.Logger(JaegerLogger(logger)),
		jaegerconfig.Metrics(jMetricsFactory))
	if err != nil {
		panic(fmt.Errorf("could not initialize jaeger tracer: %w", err))
	}

	return tracer, closer
}

// JaegerLogger implements the jaeger.Logger interface for logger.Logger
type JaegerLogger logger.Logger

func (l JaegerLogger) Error(msg string) {
	log := logger.Logger(l)
	log.Err(errors.New(msg)).Msg("jeager tracing error")
}

func (l JaegerLogger) Infof(msg string, args ...interface{}) {
	log := logger.Logger(l)
	log.Info().Msgf("%q", fmt.Sprintf(msg, args...))
}

func (l JaegerLogger) Debugf(msg string, args ...interface{}) {
	log := logger.Logger(l)
	log.Debug().Msgf("%q", fmt.Sprintf(msg, args...))
}
