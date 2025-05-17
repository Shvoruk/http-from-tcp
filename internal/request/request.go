package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Split at the first \r\n to get just the request line
	text := string(data)
	lines := strings.SplitN(text, "\r\n", 2)
	if len(lines) == 0 || lines[0] == "" {
		return nil, errors.New("missing request line")
	}

	rl, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *rl}, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid number of parts in request line")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	// Validate method: only uppercase A-Z
	for _, r := range method {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return nil, errors.New("invalid method")
		}
	}

	// Validate version: only HTTP/1.1 for now
	if version != "HTTP/1.1" {
		return nil, errors.New("unsupported HTTP version")
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, nil
}
