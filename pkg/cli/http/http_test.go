package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	h "github.com/spacelavr/dlm/pkg/cli/http"
	"github.com/stretchr/testify/assert"
)

func TestGET(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_, err := h.GET(ts.URL)
	assert.NoError(t, err)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	_, err = h.GET(ts.URL)
	assert.Error(t, err)
}
