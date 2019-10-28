package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndersonQ/go-skeleton/constants"
)

func TestRequestIDHandlerWithRequestID(t *testing.T) {
	want := "requestID"

	r := httptest.NewRequest("ignore", "/ignore", nil)
	r.Header.Set(constants.HeaderKeyRequestID, want)
	w := httptest.NewRecorder()

	h := RequestIDHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	h.ServeHTTP(w, r)

	got := w.Header().Get(constants.HeaderKeyRequestID)

	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestRequestIDHandlerWithoutRequestID(t *testing.T) {
	r := httptest.NewRequest("ignore", "/ignore", nil)
	w := httptest.NewRecorder()
	var want string

	h := RequestIDHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := r.Context().Value(reaIDContextKey).(string)

		if !ok {
			t.Fatalf("requestID is not a string! got: %s", got)
		}
		if got == "" {
			t.Fatalf("requestID from request context is empty, want a uuid as requestID")
		}

		want = got
	}))

	h.ServeHTTP(w, r)
	got := w.Header().Get(constants.HeaderKeyRequestID)

	if got == "" {
		t.Errorf("got: %s, want an uuid", got)
	}

	if got != want {
		t.Errorf("got requestID %s, want: %s", got, got)
	}
}
