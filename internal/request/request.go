package request

import (
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	// Headers     map[string]string
	// Body        []byte
}

// validation check for request line header
func (r *RequestLine) ValidateHTTP() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		//not found the start line, return nil
		return nil, b, nil
	}

	startLine := b[:idx]
	restOfMsg := b[idx+len(SEPERATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, ERROR_MALFORMED_REQUEST_LINE
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}
	if rl.ValidateHTTP() {
		return rl, restOfMsg, nil
	}
	return nil, restOfMsg, ERROR_UNSUPORTED_HTTP_VERSION

}

func RequestFromReader(reader io.Reader) (*Request, error) {

}
