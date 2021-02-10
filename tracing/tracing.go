package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// SpanFromCtx returns the non-nil span returned by opentracing.SpanFromContext
// or a span from the noop tracer if there is no span in the context.
func SpanFromContext(ctx context.Context) opentracing.Span {
	sp := opentracing.SpanFromContext(ctx)
	if sp == nil {
		return opentracing.NoopTracer{}.StartSpan("noopTracer")
	}

	return sp
}
