package request

var SEPERATOR = []byte("\r\n")

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)
