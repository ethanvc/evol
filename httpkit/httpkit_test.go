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
