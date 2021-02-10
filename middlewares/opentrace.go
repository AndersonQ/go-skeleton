package middlewares

import (
	"net/http"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/opentracing/opentracing-go"
)

// AddOpentracing adds an opentracing span to the context and finishes the span
// when handler returns.
// Use tracing.SpanFromContext to get the span from the context. It is
// technically safe to call opentracing.SpanFromContext after this middleware
// and trust the returned span is not nil. However tracing.SpanFromContext is
// safer as it'll return a disabled span if none is found in the context.
func AddOpentracing(path string, tracer opentracing.Tracer, handler http.Handler) (string, http.HandlerFunc) {
	return path, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := extractSpan(r, path, tracer)
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)

		span.SetTag("tracking_id", tracking.IDFromContext(ctx))

		err := tracer.Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(w.Header()))
		if err != nil {
			logger.FromContext(ctx).Err(err).Msg("could not inject opentracing span")
		}

		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}

func extractSpan(r *http.Request, path string, tracer opentracing.Tracer) opentracing.Span {
	carrier := opentracing.HTTPHeadersCarrier(r.Header)
	logger.FromContext(r.Context()).Debug().Msgf("r.Header: %#v", r.Header)
	logger.FromContext(r.Context()).Debug().Msgf("carrier: %#v", carrier)

	spanContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err == nil {
		logger.FromContext(r.Context()).Err(err).Msg("could not extract span")
	}

	logger.FromContext(r.Context()).Debug().Msgf("span context: %T", spanContext)
	logger.FromContext(r.Context()).Debug().Msgf("span context: %#v", spanContext)
	span := tracer.
		StartSpan(path,
			opentracing.ChildOf(spanContext))

	return span
}
