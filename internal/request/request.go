package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal("error", "error", err)
	}
	rl, err := parseRequestLine(r)
	request := Request{
		RequestLine: rl,
	}
	return &request, err
}

func parseRequestLine(request []byte) (RequestLine, error) {
	// init
	requestLine := RequestLine{}

	// split the request on "\r\n"
	lines := bytes.SplitN(request, []byte("\r\n"), 2)
	if len(lines) == 0 {
		return requestLine, fmt.Errorf("empty request")
	}

	// split the request line into its parts separated by space
	requestLineParts := strings.Split(string(lines[0]), " ")

	if n := len(requestLineParts); n != 3 {
		return requestLine, fmt.Errorf(
			"invalid number of parts in requestline: %d",
			n,
		)
	}

	requestLine.Method = requestLineParts[0]
	requestLine.RequestTarget = requestLineParts[1]
	requestLine.HttpVersion = strings.SplitAfter(requestLineParts[2], "/")[1]

	return requestLine, nil
}
