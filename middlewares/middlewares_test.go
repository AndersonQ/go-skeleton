package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AndersonQ/go-skeleton/handlers"
)

func TestJsonResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("ignore", "/ignore", nil)

	h := JsonResponse(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	h.ServeHTTP(w, r)

	got := w.Result().Header.Get(handlers.ContentType)

	if got != handlers.ContentTypeJSON {
		t.Errorf("got header %s: %s, want: %s", handlers.ContentType, got, handlers.ContentTypeJSON)
	}
}

func TestTimeoutWrapper(t *testing.T) {
	want := http.StatusServiceUnavailable
	timeout := time.Millisecond

	w := httptest.NewRecorder()
	r := httptest.NewRequest("ignore", "/ignore", nil)

	h := TimeoutWrapper(timeout)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * timeout)
	}))

	h.ServeHTTP(w, r)

	got := w.Result().StatusCode
	if got != want {
		t.Errorf("got http status: %d, want: %d", got, want)
	}
}
