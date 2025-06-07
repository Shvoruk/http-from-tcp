package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}
	line := string(data[:idx])
	colonIdx := strings.Index(line, ":")
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("invalid header line: %v", line)
	}
	name := line[:colonIdx]
	if strings.TrimSpace(name) != name {
		return 0, false, fmt.Errorf("invalid header name: %v", name)
	}

	name = strings.TrimSpace(name)
	value := strings.TrimSpace(line[colonIdx+1:])

	h[name] = value
	return idx + 2, false, nil
}
