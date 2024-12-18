package httpkit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendRequest(t *testing.T) {
	sa := NewSingleAttempt(context.Background(), http.MethodGet, "http://www.baidu.com")
	var respStr string
	err := SendRequest(sa, nil, &respStr)
	require.NoError(t, err)
}

func TestJsonResponse(t *testing.T) {
	type Project struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	sa := NewSingleAttempt(context.Background(), http.MethodGet, "https://gitlab.com/api/v4/projects")
	var resp []Project
	err := SendRequest(sa, nil, &resp)
	require.NoError(t, err)
}
