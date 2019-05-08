package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"

	"github.com/AndersonQ/go-skeleton/constants"
)

const ContentType = "Content-Type"
const ContentTypeJson = "application/json; charset=utf-8"

var version = "development"
var buildTime = "build tome not set"

// NewLivenessHandler a handler for kubernetes liveness probe
func NewLivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ContentType, ContentTypeJson)
		resp := fmt.Sprintf(
			`{"status":"Kubernetes I'm ok', no need to restart me,"version":"%s","build_time":"%s"}`,
			version,
			buildTime)

		_, _ = w.Write([]byte(resp))
	}
}

// NewLivenessHandler a handler for kubernetes readness probe
func NewReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ContentType, ContentTypeJson)
		resp := fmt.Sprintf(
			`{"status":"Kubernetes I'm ok', you can send requests to me,"version":"%s","build_time":"%s"}`,
			version,
			buildTime)

		_, _ = w.Write([]byte(resp))
	}
}

func RequestLogWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Use a clock abstraction instead of time.Now()
		startTime := time.Now()

		logger := zerolog.Ctx(r.Context()).With().
			Str(constants.LogKeyHTTPMethod, r.Method).
			Str(constants.LogKeyURLPath, r.URL.Path).
			Str(constants.LogKeyUserAgent, r.UserAgent()).
			Str(constants.LogKeyRemoteAddr, r.RemoteAddr).
			Logger()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			logger.Info().
				Fields(map[string]interface{}{
					constants.LogKeyResquestDuration: time.Now().Sub(startTime),
					constants.LogKeyHTTPStatus:       ww.Status()}).
				Msgf("%s %s", r.Method, r.RequestURI)
		}()

		next.ServeHTTP(ww, r)
	})
}
