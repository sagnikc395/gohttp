package request

import "io"

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

func RequestFromReader(reader io.Reader) (*Request, error) {
	
}
