package headers

import (
	"bytes"
	"fmt"
)

var rn = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0

	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			return read, false, nil
		}

		//check for the empty line (end of headers)
		if idx == 0 {
			read += len(rn)
			return read, true, nil
		}

		//data[read: read+idx] to get the Key:Value line
		line := data[read : read+idx]
		name, value, err := parseHeader(line)
		if err != nil {
			return 0, false, err
		}

		h[name] = value
		read += idx + len(rn)
	}
}
