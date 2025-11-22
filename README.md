changing implementation
***
# ServeLite

A lightweight HTTP server built from scratch in Go without using the builtin server.

## Features

- HTTP/1.1 support with custom parsing (no `net/http` dependency)
- State-machine based request parser
- Concurrent request handling with goroutines
- Chunked transfer encoding for streaming responses
- Zero-copy header manipulation
- Signal-based graceful shutdown (SIGINT/SIGTERM)

## Installation

```bash
git clone <github.com/srivastavcodes/ServeLite>
cd servelite
go build -o server ./cmd/server
```

## Usage

Start the server:

```bash
./server
```

The server will start on `localhost:7714` by default.

Test endpoints:

```bash
curl http://localhost:7714/                    # 200 OK
curl http://localhost:7714/yourproblem         # 400 Bad Request  
curl http://localhost:7714/myproblem           # 500 Internal Server Error
curl http://localhost:7714/httpbin/stream/5    # Chunked streaming proxy
```

## Architecture

Built with modular components:

- **Request Parser**: RFC-compliant HTTP parsing with proper error handling
- **Response Writer**: Enforced state transitions (status → headers → body)
- **Header Manager**: Case-insensitive, multi-value header support
- **Server Core**: TCP listener with per-connection goroutines

## Configuration

```bash
go run ./cmd/server
```

Change port by modifying the `port` constant in `cmd/server/main.go`:

```go
const port = 7714
```

Run tests:

```bash
go test ./...
```
