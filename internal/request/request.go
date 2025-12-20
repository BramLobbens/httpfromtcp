package request

import (
	"io"
	"log"
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
	foo, err := parseRequestLine((string(r)))
	if err != nil {
		log.Fatal("error", "error", err)
	}
	return nil, nil
}

func parseRequestLine(request string) (string, error) {
	return "", nil
}
