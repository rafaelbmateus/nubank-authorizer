package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"nubank/authorizer/database/memory"
	"nubank/authorizer/handler"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()
var log = &zerolog.Logger{}

func TestRoutes(t *testing.T) {
	ts := httptest.NewServer(handler.NewAuthorizer(ctx, &memory.Memory{}, log).Server())
	defer ts.Close()
	resp, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHome(t *testing.T) {
	route := gin.Default()
	route.GET("/", handler.Home)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
