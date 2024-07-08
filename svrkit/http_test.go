package svrkit

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_BasicHttp(t *testing.T) {
	engine := gin.New()
	kit := NewGinKit()
	engine.POST("/", kit.Handlers(func(c context.Context, req any, nexter Nexter) (any, error) {
		return nil, nil
	}))
	httpReq := httptest.NewRequest(http.MethodPost, "http://www.xx.com/", strings.NewReader("test"))
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httpReq)
	require.Equal(t, http.StatusOK, w.Code)
}
