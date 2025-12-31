package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) (string, bool) {
	val, ok := h[strings.ToLower(key)]
	return val, ok
}

var isTokenChar = [256]bool{
	'!': true, '#': true, '$': true, '%': true, '&': true,
	'\'': true, '*': true, '+': true, '-': true, '.': true,
	'^': true, '_': true, '`': true, '|': true, '~': true,
}

func init() {
	for c := byte('a'); c <= 'z'; c++ {
		isTokenChar[c] = true
	}
	for c := byte('A'); c <= 'Z'; c++ {
		isTokenChar[c] = true
	}
	for c := byte('0'); c <= '9'; c++ {
		isTokenChar[c] = true
	}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlf := []byte("\r\n")

	startsWithCrlf := bytes.HasPrefix(data, crlf)
	// check for end of headers if at start
	if startsWithCrlf {
		// empty line - end of content
		return len(crlf), true, nil
	}

	fieldLineEnd := bytes.LastIndex(data, crlf)
	if fieldLineEnd == -1 {
		// no complete field-line in data
		return 0, false, nil
	}

	fieldLines := bytes.Split(data, crlf)
	numBytesParsed := 0

	for _, line := range fieldLines {
		if len(line) == 0 {
			break
		}

		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			// incomplete or malformed data - retry
			return numBytesParsed, false, nil
		}

		if s := strings.TrimRight(string(parts[0]), " "); len(s) < len(parts[0]) {
			return numBytesParsed, false, fmt.Errorf("invalid formatting in field-name")
		}

		key := strings.ToLower(string(bytes.TrimSpace(parts[0])))
		val := string(bytes.TrimSpace(parts[1]))

		if !isFieldNameValid(key) {
			return numBytesParsed, false, fmt.Errorf("invalid token in field-name: %s", key)
		}

		if keyVal, ok := h[key]; ok {
			// multiple values in key concatenate into a single string, separated by a comma
			newVal := strings.Join(append(strings.Split(keyVal, ", "), val), ", ")
			h[key] = newVal
		} else {
			h[key] = val
		}

		numBytesParsed += len(line) + len(crlf)
	}

	done = bytes.Index(data, []byte("\r\n\r\n")) != -1
	return numBytesParsed, done, nil
}

func isFieldNameValid(fieldName string) bool {
	isTokenInvalid := func(r rune) bool {
		return !isTokenChar[r]
	}
	return len(fieldName) >= 1 && !strings.ContainsFunc(fieldName, isTokenInvalid)
}
