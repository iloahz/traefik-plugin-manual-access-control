package plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNoCookie(t *testing.T) {
	cfg := CreateConfig()
	cfg.Server = "http://172.16.11.24:9502"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "tpmac")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	t.Log(recorder.Result().StatusCode)
	t.Log(recorder.Result().Cookies())
}
