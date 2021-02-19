package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"

	"github.com/AndersonQ/go-skeleton/tracing"
)

func NewOpentracing() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sp := tracing.SpanFromContext(ctx)
		sp.LogFields(
			log.String("some_id", uuid.New().String()),
			log.String("some_key", "some_value"),
			log.Int("some_number", 1500),
			log.Event("some_event"),
			log.Message("some message"))

		if rand.Int()%2 == 0 {
			sp.LogFields(log.Error(errors.New("bad coin flip")))
			ext.Error.Set(sp, true)
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte(`I'm a tea pot`))
			return
		}

		url := r.URL.Query().Get("url")
		if url != "" {
			if err := externalCall(ctx, url); err != nil {
				logger.FromContext(ctx).Err(err).Msg("external call failed")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"error":%q,"tracking_id":"%s"}`,
					err, tracking.IDFromContext(ctx))))
				return
			}
		}

		headers, _ := json.Marshal(w.Header())
		_, _ = w.Write([]byte(
			fmt.Sprintf(`{"tracking_id":"%s","url":"%s","headers":%q}`,
				tracking.IDFromContext(ctx),
				url,
				headers)))
	})
}

func externalCall(ctx context.Context, url string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.FromContext(ctx).Err(err).Msgf("create request to %q", url)
	}

	req.Header.Set("X-Tracking-Id", tracking.IDFromContext(ctx))
	sp := tracing.SpanFromContext(ctx)
	err = sp.Tracer().Inject(
		sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		logger.FromContext(ctx).Err(err).Msg("could not inject opentracing span")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call %q: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%q replied with status %s", url, resp.Status)
	}

	return nil
}
