package middlewares

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var middleware http.Handler

type FakeResponse struct {
	t       *testing.T
	headers http.Header
	body    []byte
	status  int
}

func New(t *testing.T) *FakeResponse {
	return &FakeResponse{
		t:       t,
		headers: make(http.Header),
	}
}

func (r *FakeResponse) Header() http.Header {
	return r.headers
}

func (r *FakeResponse) Write(body []byte) (int, error) {
	r.body = body
	return len(body), nil
}

func (r *FakeResponse) WriteHeader(status int) {
	r.status = status
}

func (r *FakeResponse) Assert(status int, body string) {
	if r.status != status {
		r.t.Errorf("expected status %+v to equal %+v", r.status, status)
	}
	if string(r.body) != body {
		r.t.Errorf("expected body %+v to equal %+v", string(r.body), body)
	}
}

func TestNewStatusMiddleware(t *testing.T) {
	middleware = NewStatusMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	assert.NotEmpty(t, middleware)
}

func TestHeader(t *testing.T) {

}

func TestWrite(t *testing.T) {

}

func TestWriteHeader(t *testing.T) {

}

func TestStatusCode(t *testing.T) {

}
