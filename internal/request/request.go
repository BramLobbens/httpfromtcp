package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

type requestState int

const (
	Initialized requestState = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	State       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request{}
	lines, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal("error", "error", err)
	}
	err = req.parseRequestLine(lines)
	return &req, err
}

// func (r *Request) parse(data []byte) (int, error) {

// }

func (req *Request) parseRequestLine(request []byte) error {
	// split the request on "\r\n"
	lines := bytes.SplitN(request, []byte("\r\n"), 2)
	if len(lines) == 0 {
		return fmt.Errorf("empty request")
	}

	// split the request line into its parts separated by space
	parts := strings.Split(string(lines[0]), " ")

	if n := len(parts); n != 3 {
		return fmt.Errorf(
			"invalid number of parts in requestline: %d",
			n,
		)
	}

	if err := req.RequestLine.parseParts(parts); err != nil {
		return err
	}
	return nil
}

func (rl *RequestLine) parseParts(parts []string) error {
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
