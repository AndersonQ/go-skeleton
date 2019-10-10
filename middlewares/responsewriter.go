package middlewares

import "net/http"

// statusResponseWriter wraps a http.ResponseWriter to expose the status code through statusResponseWriter.statusCode
type statusResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
}

func NewStatusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := statusResponseWriter{w: w}

		next.ServeHTTP(&ww, r)
	})
}

func (s *statusResponseWriter) Header() http.Header {
	return s.w.Header()
}

func (s *statusResponseWriter) Write(bs []byte) (int, error) {
	return s.w.Write(bs)
}

func (s *statusResponseWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.w.WriteHeader(statusCode)
}

func (s *statusResponseWriter) StatusCode() int {
	return s.statusCode
}
