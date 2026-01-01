package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/BramLobbens/httpfromtcp/internal/headers"
)

const bufferSize = 8

type requestState int

const (
	requestStateInitialized    requestState = iota // 0
	requestStateParsingHeaders                     // 1
	requestStateParsingBody                        // 2
	requestStateDone                               // 3
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}
	buffer := make([]byte, bufferSize)
	readToIndex := 0

	for req.state != requestStateDone {
		if readToIndex >= len(buffer) {
			// Grow buffer - double the capacity
			newBuffer := make([]byte, cap(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		fmt.Fprintf(os.Stderr, "DEBUG: readToIndex=%d, len(buffer)=%d, cap(buffer)=%d, buffer[readToIndex:] len=%d\n",
			readToIndex, len(buffer), cap(buffer), len(buffer[readToIndex:]))

		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if numBytesRead > 0 {
			readToIndex += numBytesRead
		}
		if readToIndex > 0 {
			// always process the n > 0 bytes returned before considering the error err
			numBytesParsed, err := req.parse(buffer[:readToIndex])
			if numBytesParsed > 0 {
				// Copy remaining data to the start
				remainingBytes := readToIndex - numBytesParsed
				copy(buffer, buffer[numBytesParsed:readToIndex])
				readToIndex = remainingBytes

				// Clear the rest of the buffer to avoid confusion
				for i := readToIndex; i < len(buffer); i++ {
					buffer[i] = 0
				}
			}
			if err != nil {
				return req, err
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != requestStateDone {
					return nil, fmt.Errorf("unexpected end of input")
				}
				break
			}
			return nil, err
		}
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	//for r.State != requestStateDone {
	numBytesParsed, err := r.parseSingle(data[totalBytesParsed:])
	if err != nil {
		return numBytesParsed, err
	}
	totalBytesParsed += numBytesParsed
	//}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		numOfBytes, err := r.parseRequestLine(data)
		if err != nil {
			return numOfBytes, err
		}
		if numOfBytes != 0 {
			// Finalised request line - set to parse headers next
			r.state = requestStateParsingHeaders
		}
		return numOfBytes, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateParsingBody
			return n, nil
		}
		return n, nil
	case requestStateParsingBody:
		n, err, done := r.parseBody(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateDone
		}
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}
	return 0, fmt.Errorf("error: unknown state")
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	crlf := []byte("\r\n")
	// split the request on the first occurrence of "\r\n" if found
	firstLine, _, found := bytes.Cut(data, crlf)
	numOfBytes := len(firstLine) + len(crlf)
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

	if err := r.RequestLine.validateAndSetRequestLineParts(parts); err != nil {
		return numOfBytes, err
	}

	return numOfBytes, nil
}

func (rl *RequestLine) validateAndSetRequestLineParts(parts []string) error {
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

func (r *Request) parseBody(data []byte) (int, error, bool) {
	contentLength, found := r.Headers.Get("content-length")
	if !found || contentLength == "0" {
		// no body present
		return 0, nil, true
	}

	dataLength := len(data)
	if cl, _ := strconv.Atoi(contentLength); dataLength > cl {
		return 0, fmt.Errorf("more data received than specified in content-length header"), false
	} else if dataLength == cl {
		// body complete - attention data is a slice and we're setting Body []byte, so copy the data!
		r.Body = make([]byte, len(data))
		copy(r.Body, data)
		return len(r.Body), nil, true
	}

	return 0, nil, false
}
