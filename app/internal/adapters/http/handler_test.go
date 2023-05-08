package http

import (
	"github.com/shamank/wb-l0/app/internal/service"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler(&service.Service{})

	require.IsType(t, &Handler{}, h)

}

func TestHandler_InitAPI(t *testing.T) {
	h := NewHandler(&service.Service{})
	router := h.InitAPI()

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
