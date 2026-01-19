package request

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestFromReader(t *testing.T) {
	request := "POST /order/123?image=69 HTTP/1.1\r\nHost: localhost:8080\r\nUser-Agent: Go-HTTP-Server\r\n" +
		"Accept: */*\r\n\r\nhello world"
	reader := strings.NewReader(request)

	req, err := RequestFromReader(bufio.NewReader(reader))
	require.NoError(t, err)

	require.Equal(t, "POST", req.Method)
	require.Equal(t, "/order/123", req.URL)
	require.Equal(t, "69", req.QueryParams["image"])
	require.Equal(t, "Go-HTTP-Server", req.Headers.Headers["user-agent"])
}
