package request

import (
	"bytes"
	"fmt"
	"io"

	"github.com/sagnikc395/gohttp/internal/headers"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       parserState
}

func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

var SEPERATOR = []byte("\r\n")

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type parserState string

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		//not found the start line, return nil
		return nil, 0, nil
	}

	startLine := b[:idx]
	//restOfMsg := b[idx+len(SEPERATOR):]
	read := idx + len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	//Splitting the "HTTP/1.1" into ["HTTP","1.1"]
	protoParts := bytes.Split(parts[2], []byte("/"))
	if len(protoParts) != 2 || string(protoParts[0]) != "HTTP" {
		return nil, 0, ERROR_UNSUPORTED_HTTP_VERSION
	}

	return &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(protoParts[1]),
	}, read, nil

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
		if n > 0 {
			bufIdx += n
			readN, parseErr := request.parse(buf[:bufIdx])
			if parseErr != nil {
				return nil, parseErr
			}

			//sliding over the buffer and removing what waas consumed
			if readN > 0 {
				copy(buf, buf[readN:bufIdx])
				bufIdx -= readN
			}
		}
		if err != nil {
			if err == io.EOF && !request.done() {
				return nil, fmt.Errorf("connection closed before request finished")
			}
			return nil, err
		}
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {

	readTotal := 0

	for {
		switch r.State {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(data[readTotal:])
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				return readTotal, nil
			}
			r.RequestLine = *rl
			readTotal += n
			r.State = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(data[readTotal:])
			if err != nil {
				r.State = StateError
				return 0, err
			}
			readTotal += n
			if done {
				r.State = StateDone
				return readTotal, nil
			}
			//not enough data for full headers yet
			return readTotal, nil
		case StateDone:
			return readTotal, nil

		default:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		}
	}

}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}
