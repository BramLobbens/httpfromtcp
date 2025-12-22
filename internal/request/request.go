package request

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"
)

type requestState int

const crlf = "\r\n"
const bufferSize = 8

const (
	Initialized requestState = iota // 0
	Done                            // 1
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
	req.State = Initialized
	buffer := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	for req.State != Done {
		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if numBytesRead > 0 {
			// always process the n > 0 bytes returned before considering the error err
			numBytesParsed, err := req.parse(buffer[:readToIndex])
			if numBytesParsed > 0 {
				tmpSlice := make([]byte, readToIndex, readToIndex)
				copy(tmpSlice, buffer)
				readToIndex -= numBytesParsed
			}
			if err != nil {
				return &req, err
			}
		}
		if err != nil {
			if err == io.EOF {
				req.State = Done // end of stream or data
			}
			return &req, err
		}

		if readToIndex == len(buffer) {
			buffer = slices.Grow(buffer, bufferSize*2) // grow the buffer capacity
			buffer = buffer[:cap(buffer)]              // set length to full capacity (reader only reads len(p)
		}

		readToIndex += numBytesRead
	}
	return &req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		numOfBytes, err := r.parseRequestLine(data)
		if err != nil {
			return numOfBytes, err
		}
		if numOfBytes != 0 {
			r.State = Done
		}
		return numOfBytes, nil
	case Done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}
	return 0, fmt.Errorf("error: unknown state")
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	// split the request on the first occurrence of "\r\n" if found
	firstLine, _, found := bytes.Cut(data, []byte(crlf))
	numOfBytes := len(firstLine)
	// terminal \r\n not yet in chunk return with initialized state
	if !found {
		return 0, nil
	}
	// split the request line into its parts separated by space
	parts := strings.Fields(string(firstLine)) // Fields is equal to Split on " "
	if n := len(parts); n != 3 {
		return numOfBytes, fmt.Errorf(
			"invalid number of parts in requestline: %d",
			n,
		)
	}

	if err := r.RequestLine.parseParts(parts); err != nil {
		return numOfBytes, err
	}

	return numOfBytes, nil
}

func (rl *RequestLine) parseParts(parts []string) error {
	// Method
	method := parts[0]
	if method != strings.ToUpper(method) {
		return fmt.Errorf("method token is case-sensitive; standardized methods are defined in all-uppercase")
	}
	rl.Method = method

	// Target
	if strings.IndexRune(parts[1], '/') != 0 {
		return fmt.Errorf("the client MUST send \"/\" as the path within the origin-form of request-target")
	}
	rl.RequestTarget = parts[1]

	// HttpVersion
	version := strings.SplitAfter(parts[2], "/")[1]
	if version != "1.1" {
		return fmt.Errorf("only version 1.1 is supported; provided: '%v'", version)
	}
	rl.HttpVersion = version

	return nil
}
