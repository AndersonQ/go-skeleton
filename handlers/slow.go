package handlers

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func NewSlowHandler(delay time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())

		time.Sleep(delay)

		_, err := w.Write([]byte(`{}`))
		if err != nil {
			logger.Error().Err(err).Msg("slow handler")
		}
	}
}
