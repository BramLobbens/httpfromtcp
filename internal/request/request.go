package request

import (
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
	requestLine := RequestLine{
		HttpVersion:   "",
		RequestTarget: "",
		Method:        "",
	}

	// split the request on "\r\n"
	parts := strings.Split(string(request), "\r\n")

	// split the request line into its parts separated by space
	requestLineParts := strings.Split(parts[0], " ")

	if n := len(requestLineParts); n != 3 {
		return requestLine, fmt.Errorf(
			"invalid number of parts in requestline: %d",
			n,
		)
	}

	for _, p := range requestLineParts {
		if p == "GET" || p == "POST" {
			requestLine.Method = p
		} else if strings.IndexByte(p, '/') == 0 {
			requestLine.RequestTarget = p
		} else if strings.Contains(p, "HTTP/") {
			requestLine.HttpVersion = strings.SplitAfter(p, "/")[1]
		}
	}
	return requestLine, nil
}
