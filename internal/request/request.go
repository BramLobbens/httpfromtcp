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
	rl := RequestLine{}

	// split the request on "\r\n"
	lines := bytes.SplitN(request, []byte("\r\n"), 2)
	if len(lines) == 0 {
		return rl, fmt.Errorf("empty request")
	}

	// split the request line into its parts separated by space
	parts := strings.Split(string(lines[0]), " ")

	if n := len(parts); n != 3 {
		return rl, fmt.Errorf(
			"invalid number of parts in requestline: %d",
			n,
		)
	}

	if err := parseRequestLineParts(&rl, parts); err != nil {
		return rl, err
	}
	return rl, nil
}

func parseRequestLineParts(rl *RequestLine, parts []string) error {
	method := parts[0]
	if method != strings.ToUpper(method) {
		return fmt.Errorf("invalid casing in method name")
	}
	rl.Method = method

	rl.RequestTarget = parts[1]

	version := strings.SplitAfter(parts[2], "/")[1]
	if version != "1.1" {
		return fmt.Errorf("only version 1.1 is supported; provided: '%v'", version)
	}
	rl.HttpVersion = version

	return nil
}
