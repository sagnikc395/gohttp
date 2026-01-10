package request

import "fmt"

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var SEPERATOR = "\r\n"
var ERROR_UNSUPORTED_HTTP_VERSION = fmt.Errorf("unsupported http-version")
