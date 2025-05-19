package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInit requestState = iota
	requestStateDone
)

const crlf = "\r\n"
const buffSize = 8

func RequestFromReader(r io.Reader) (*Request, error) {

	buff := make([]byte, buffSize)
	readToIndex := 0

	req := &Request{
		state: requestStateInit,
	}
	for req.state != requestStateDone {
		if readToIndex >= len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		nbr, err := r.Read(buff[readToIndex:])
		if err != nil {
			if err == io.EOF {
				req.state = requestStateDone
				break
			}
			return nil, err
		}
		readToIndex += nbr

		nbp, err := req.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[nbp:])
		readToIndex -= nbp
	}
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {

	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	rl, err := requestLineFromString(string(data[:idx]))
	if err != nil {
		return nil, 0, err
	}
	return rl, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line (length != 3): %s", str)
	}

	method := parts[0]
	// Validate method
	for _, r := range method {
		if !unicode.IsUpper(r) || unicode.IsLetter(r) == false {
			return nil, fmt.Errorf("invalid request method: %s", method)
		}
	}

	target := parts[1]

	version := parts[2]
	// Validate version
	if version != "HTTP/1.1" {
		return nil, fmt.Errorf("invalid HTTP version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInit:
		rl, nb, err := parseRequestLine(data)
		if err != nil {
			// Something went wrong
			return 0, err
		}
		if nb == 0 {
			// Need more data
			return 0, nil
		}
		r.RequestLine = *rl
		r.state = requestStateDone
		return nb, nil
	case requestStateDone:
		return 0, fmt.Errorf("request already parsed")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
