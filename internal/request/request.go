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
	State parserState
}

func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

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

	parts = strings.Split(parts[2], "/")

	if len(parts) != 3 {
		return nil, restOfMsg, ERROR_MALFORMED_REQUEST_LINE
	}

	//HTTP validation , should be 1.1 only
	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMsg, ERROR_MALFORMED_REQUEST_LINE
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return rl, restOfMsg, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	//TODO: buffer could get overrun ...
	// a header that exceeds 1k or a body
	buf := make([]byte, 1024)
	//n := 0
	bufIdx := 0
	for !request.done() {
		//simulate reading slowly over time , dont want to read all at once
		// the body parsing doesnt need to happen right away

		n, err := reader.Read(buf[bufIdx:])
		if err != nil {
			return nil, err
		}

		//new posn of bufIdx
		bufIdx += n

		readN, err := request.parse(buf[:bufIdx+n])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufIdx])
		bufIdx -= readN
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {

}

func (r *Request) done() bool {
	return r.State == StateDone
}
