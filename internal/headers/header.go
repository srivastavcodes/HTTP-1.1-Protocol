package headers

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	Headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		Headers: make(map[string]string),
	}
}

func (h *Headers) Get(name string) string {
	val := h.Headers[strings.ToLower(name)]
	return val
}

func (h *Headers) Set(name string, value string) {
	name = strings.ToLower(name)

	if v, ok := h.Headers[name]; ok {
		h.Headers[name] = fmt.Sprintf("%s, %s", v, value)
	} else {
		h.Headers[name] = value
	}
}

func (h *Headers) Replace(name string, value string) {
	h.Headers[name] = strings.ToLower(value)
}

func (h *Headers) Exists(name string) bool {
	_, ok := h.Headers[strings.ToLower(name)]
	return ok
}

func (h *Headers) ParseHeader(reader *bufio.Reader) error {
	for {
		// reads one header line from the buffer
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return fmt.Errorf("could not parse header %s: %w", line, err)
		}
		// removes \r\n and checks end of Headers
		line = bytes.TrimSpace(line)
		if bytes.Equal(line, []byte("")) {
			return nil
		}
		// parsesHeader
		name, value, err := h.parseHeader(line)
		if err != nil {
			return fmt.Errorf("could not parse header %s: %w", line, err)
		}
		h.Set(name, value)
	}
}

func (h *Headers) parseHeader(line []byte) (string, string, error) {
	header := bytes.SplitN(line, []byte(":"), 2)
	if len(header) != 2 {
		return "", "", fmt.Errorf("malformed feild line: %s", line)
	}
	var (
		name  = header[0]
		value = bytes.TrimSpace(header[1])
	)
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed header name: %s", line)
	}
	if !isValidToken(name) {
		return "", "", fmt.Errorf("malformed header name: %s", line)
	}
	return string(name), string(value), nil
}

func isValidToken(chars []byte) bool {
	for _, ch := range chars {
		var result bool
		if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9' {
			result = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			result = true
		}
		if !result {
			return false
		}
	}
	return true
}
