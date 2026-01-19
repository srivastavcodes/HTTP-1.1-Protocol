package response

import (
	"fmt"
	"io"
	"net/http"
	"sl/internal/headers"
	"sl/internal/server/html"
	"strconv"
	"time"
)

type ResponseWriter interface {
	Header() *headers.Headers
	Write([]byte) (n int, err error)
}

type HTTPResponse struct {
	Headers *headers.Headers
	Code    int
	Text    string

	io.Writer

	Body string
}

// NewHTTPResponse creates a response with default values.
func NewHTTPResponse() *HTTPResponse {
	return &HTTPResponse{
		Code:    http.StatusOK,
		Text:    "OK",
		Headers: headers.NewHeaders(),
		Body:    "",
	}
}

// Header returns the underlying header map.
func (r *HTTPResponse) Header() *headers.Headers {
	return r.Headers
}

// SetStatus sets the Code and Text in HTTPResponse
func (r *HTTPResponse) SetStatus(code int, text string) {
	r.Code = code
	r.Text = text
}

// SetHeader sets the headers in HTTPResponse.Headers
func (r *HTTPResponse) SetHeader(key, value string) {
	r.Headers.Set(key, value)
}

// SetBody sets the Body in HTTPResponse.Body and also write the
// Content-Length header.
func (r *HTTPResponse) SetBody(body string) {
	r.Headers.Replace("Content-Length", strconv.Itoa(len(body)))
	r.Body = body
}

/*
	HTTP/1.1 200 OK\r\n
	Header-Name: Header-Value\r\n
	Content-Type: application/json\r\n
	Content-Length: 200\r\n
	\r\n
	Response Body
*/

// WriteResponse converts HTTPResponse to raw bytes to send to the
// client.
func (r *HTTPResponse) WriteResponse() []byte {
	response := make([]byte, 0)
	response = fmt.Appendf(response, "HTTP/1.1 %d %s\r\n", r.Code, r.Text)

	for key, value := range r.Header().Headers {
		response = fmt.Appendf(response, "%s: %s\r\n", key, value)
	}
	response = append(response, []byte("\r\n")...)
	response = append(response, []byte(r.Body)...)
	return response
}

func SayHello() *HTTPResponse {
	res := NewHTTPResponse()

	res.SetHeader("Content-Type", "text/html; charset=utf-8")
	res.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	res.SetHeader("Server", "GrugHTTPServer/1.1")
	res.SetHeader("Connection", "close")
	res.SetHeader("Date", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))

	// Claude's HTML skills says hello
	res.Body = html.SayHelloHtml
	return res
}
