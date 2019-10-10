package middlewares

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/AndersonQ/go-skeleton/constants"
	"github.com/AndersonQ/go-skeleton/handlers"
)

// JsonResponse sets the response content-type header to application/json; charset=utf-8
// it will overwrite the content-type set by any inner middleware
func JsonResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// overwrite the content-type set by any inner middleware
		defer func() {
			w.Header().Set(handlers.ContentType, handlers.ContentTypeJSON)
		}()

		next.ServeHTTP(w, r)
	})
}

// RequestLogWrapper adds a zerolog.Loggger to the request context, set http specific log
// fields and at the end of the request logs the request details and its duration in ms
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

		ww := statusResponseWriter{w: w}

		defer func() {
			logger.Info().
				Fields(map[string]interface{}{
					constants.LogKeyResquestDuration: time.Now().Sub(startTime),
					constants.LogKeyHTTPStatus:       ww.statusCode}).
				Msgf("%s %s", r.Method, r.RequestURI)
		}()

		next.ServeHTTP(&ww, r)
	})
}

func TimeoutWrapper(timeout time.Duration) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.TimeoutHandler(handler, timeout, "")
	}
}
