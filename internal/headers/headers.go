package headers

import "bytes"

var rn = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data, rn)
		if idx == -1 {
			break
		}
	}
	return read, done, nil
}
