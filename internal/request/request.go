package request

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"sl/internal/headers"
	"strconv"
	"strings"
)

// HTTPRequest represent a parsed request.
type HTTPRequest struct {
	Method      string            // Type of request - GET, POST, etc.
	URL         string            // The requested path. /order/123
	RawURL      string            // Complete url
	QueryParams map[string]string // ../requestedUrl?image=69
	Version     string            // HTTP version - HTTP/1.1, HTTP/1.0
	Headers     *headers.Headers  // HTTP headers - Content-Type, Content-Length
	Body        string            // Body of the request
}

func RequestFromReader(reader *bufio.Reader) (*HTTPRequest, error) {
	req := HTTPRequest{
		QueryParams: make(map[string]string),
		Headers:     headers.NewHeaders(),
	}
	err := req.ParseRequestLine(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing request line: %w", err)
	}
	err = req.Headers.ParseHeader(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing request header: %w", err)
	}
	// todo -> parse body and write tests
	return &req, nil
}

// Request-Line "GET /order/123?image=69 HTTP/1.1"

// ParseRequestLine parses the first line of the entire HTTP request.
func (r *HTTPRequest) ParseRequestLine(reader *bufio.Reader) error {
	rl, err := reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("error reading request-line: %w", err)
	}
	rlSlice := bytes.Fields(bytes.TrimSpace(rl))
	if len(rlSlice) != 3 {
		return fmt.Errorf("invalid request line. expected=3 | got=%d", len(rlSlice))
	}
	r.Method = string(bytes.ToUpper(rlSlice[0]))

	r.RawURL = string(rlSlice[1])
	r.Version = string(rlSlice[2])

	parsedUrl, err := url.Parse(r.RawURL)
	if err != nil {
		return fmt.Errorf("failed to parse url %s: %w", r.RawURL, err)
	}
	r.URL = parsedUrl.Path

	for key, val := range parsedUrl.Query() {
		if len(val) > 0 {
			r.QueryParams[key] = val[0]
		}
	}
	return nil
}

// GetHeader retrieves the header value.
func (r *HTTPRequest) GetHeader(name string) string {
	return r.Headers.Get(name)
}

// AddHeader adds a header to the request.
func (r *HTTPRequest) AddHeader(name string, value string) {
	r.Headers.Set(name, value)
}

// HasHeader checks if the header exists.
func (r *HTTPRequest) HasHeader(name string) bool {
	return r.Headers.Exists(name)
}

// GetQueryParam retrieves a query parameter value.
func (r *HTTPRequest) GetQueryParam(name string) string {
	return r.QueryParams[name]
}

// HasQueryParam checks if the query parameter exists.
func (r *HTTPRequest) HasQueryParam(name string) bool {
	_, ok := r.QueryParams[name]
	return ok
}

// IsMethod checks if the request method matches.
func (r *HTTPRequest) IsMethod(method string) bool {
	return strings.EqualFold(r.Method, method)
}

// GetContentLength returns content length from headers.
func (r *HTTPRequest) GetContentLength() int {
	if cl := r.GetHeader("content-length"); cl != "" {
		if cl, err := strconv.Atoi(cl); err == nil {
			return cl
		}
	}
	return 0
}
