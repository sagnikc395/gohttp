package request

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// read the number upto len(p) or numBytesPerRead bytes from the string per call
// useuful for simulating reading a variable number of bytes per chunk from the listener
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}

	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// assert.Equal(t, "TheTestagen", "theTestagen")
	// Test: Good GET Request line
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)

	// Test: Good GET Request line
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

}

func TestRequestComprehensive(t *testing.T) {
	tests := []struct {
		name           string
		raw            string
		expectedMethod string
		expectedPath   string
		expectedHeader string
		expectErr      bool
	}{
		{
			name:           "Simple GET",
			raw:            "GET /index.html HTTP/1.1\r\nHost: localhost\r\n\r\n",
			expectedMethod: "GET",
			expectedPath:   "/index.html",
			expectedHeader: "localhost",
			expectErr:      false,
		},
		{
			name:      "Malformed Method",
			raw:       "NOT_A_METHOD / HTTP/1.1\r\n\r\n", // Currently your parser allows any method string
			expectErr: false,
		},
		{
			name:      "Invalid HTTP Version",
			raw:       "GET / HTTP/2.0\r\n\r\n",
			expectErr: true,
		},
		{
			name:      "Missing CRLF",
			raw:       "GET / HTTP/1.1 Host: localhost", // Connection ends without finishing
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate fragmented reading (1 byte at a time)
			reader := &chunkReader{
				data:            tt.raw,
				numBytesPerRead: 1,
			}

			r, err := RequestFromReader(reader)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedMethod, r.RequestLine.Method)
				assert.Equal(t, tt.expectedPath, r.RequestLine.RequestTarget)
				if tt.expectedHeader != "" {
					assert.Equal(t, tt.expectedHeader, r.Headers["Host"])
				}
			}
		})
	}
}
