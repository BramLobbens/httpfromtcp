package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

type requestState int

const _bufferSize = 8

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
	buffer := make([]byte, _bufferSize)
	var bufferCopy []byte
	readToIndex := 0
	for {
		if req.State == Done {
			return &req, nil
		}
		i, err := io.ReadAtLeast(reader, buffer, _bufferSize)
		if err != nil {
			if err == io.EOF {
				req.State = Done
				break
			}
			log.Fatal(err)
		}
		readToIndex += i // i bytes read from reader

		// to fix logic regarding copying to data out of the buffer
		// and increasing the slice to be parsed with the copied read data
		if i == _bufferSize { // buffer fully read
			bufferCopy = make([]byte, len(buffer)*2)
			copy(bufferCopy, buffer)
		}
		data := bufferCopy[:readToIndex]
		n, err := req.parse(data)
		if err != nil {
			log.Fatal(err)
		}
		if n != 0 { // succesfully parsed line
			copy(bufferCopy, buffer) // Remove the data that was parsed successfully
			readToIndex -= n         // n bytes parsed in buffer, i - n left
		}
	}
	return &req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == Initialized {
		numOfBytes, err := r.parseRequestLine(data)
		if err != nil {
			log.Fatal(err)
		}
		if numOfBytes != 0 {
			r.State = Done
		}
		return numOfBytes, nil
	} else if r.State == Done {
		fmt.Errorf("error: trying to read data in a done state")
	}
	fmt.Errorf("error: unknown state")
	return 0, nil
}

func (req *Request) parseRequestLine(data []byte) (int, error) {
	// split the request on the first occurrence of "\r\n" if found
	firstLine, _, found := bytes.Cut(data, []byte("\r\n"))
	numOfBytes := len(firstLine)
	// terminal \r\n not yet in chunk return with initialized state
	if !found {
		return 0, nil
	}
	// split the request line into its parts separated by space
	parts := strings.Fields(string(data)) // Fields is equal to Split on " "
	if n := len(parts); n != 3 {
		return numOfBytes, fmt.Errorf(
			"invalid number of parts in requestline: %d",
			n,
		)
	}

	if err := req.RequestLine.parseParts(parts); err != nil {
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
