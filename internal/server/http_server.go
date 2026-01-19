package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sl/internal/request"
	"sl/internal/response"
	"sl/internal/server/html"
	"time"
)

type HTTPServer struct {
	Host    string
	Port    string
	log     *log.Logger
	closed  bool
	handler Handler
}

type Handler func(w response.ResponseWriter, r *request.HTTPRequest)

func NewHttpServer(port string, handler Handler) *HTTPServer {
	logger := log.New(os.Stdout, "HTTP :: ", log.LstdFlags|log.Lmsgprefix)
	return &HTTPServer{
		Host:    "127.0.0.1",
		Port:    port,
		log:     logger,
		handler: handler,
	}
}

// ServeHTTP starts the HTTP server and listens for incoming connections. It
// accepts connections in a loop and spawns a goroutine to handle each one.
//
// Returns an error if the server fails to start listening on the configured
// address.
func (s *HTTPServer) ServeHTTP() error {
	address := fmt.Sprintf("%s:%s", s.Host, s.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on port=%s. err=%s", s.Port, err)
	}
	defer listener.Close()
	s.logf("Server listening on port=%s", s.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logf("error accepting connection: %s\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection processes a single HTTP client connection. It reads the
// HTTP request line and headers, routes the request based on the path, and
// sends back an appropriate HTTP response before closing the connection.
func (s *HTTPServer) handleConnection(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	s.logf("new connection from: %s\n", conn.RemoteAddr().String())
	_ = conn.SetReadDeadline(time.Now().Add(15 * time.Second))

	writer := response.NewHTTPResponse()

	// reader contains the http request
	reader := bufio.NewReader(conn)

	req, err := request.RequestFromReader(reader)
	if err != nil {
		s.logf("error parsing request: %s\n", err)
		return
	}
	s.handler(writer, req)
}

// TODO: enhance the error handling and html with dynamic data

// sendErrorResponse sends an error response back to the client with a html
// body containing error details.
func (s *HTTPServer) sendErrorResponse(conn net.Conn, code int, text string) {
	res := response.NewHTTPResponse()

	res.SetStatus(code, text)

	res.SetHeader("Content-Type", "text/html; charset=utf-8")
	res.SetHeader("Connection", "close")

	res.SetBody(html.BadRequestHTML)

	conn.Write(res.WriteResponse())
}

// logf logs defensively
func (s *HTTPServer) logf(format string, args ...any) {
	if s.log != nil {
		s.log.Printf(format, args...)
	} else {
		logger := log.New(os.Stdout, "HTTP :: ", log.LstdFlags|log.Lmsgprefix)
		logger.Printf(format, args...)
	}
}
