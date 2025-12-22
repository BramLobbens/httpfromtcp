package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"
const fieldSep = ":"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Check fieldline is complete
	fieldLine, _, found := bytes.Cut(data, []byte(crlf))
	if !found {
		return 0, false, nil
	}

	// newline is beginning of data, finished parsing
	if len(fieldLine) == 0 {
		return 0, true, nil
	}

	// Parse into field-name and field-value by separator
	name, value, found := bytes.Cut(fieldLine, []byte(fieldSep))
	if !found {
		return 0, false, fmt.Errorf("invalid formatting in field-line")
	}

	// Field-name must not have a space between the name and the separator
	if s := strings.TrimRight(string(name), " "); len(s) < len(name) {
		return 0, false, fmt.Errorf("invalid formatting in field-name")
	}

	// Trim value and add to map if not exists
	key := string(name)
	val := string(bytes.TrimSpace(value))

	if _, ok := h[key]; ok {
		return 0, false, fmt.Errorf("key-value pair already exists")
	}

	h[key] = val
	// bytes consumed, silly work-around, since we took a different approach to handling the crlf?
	return len(fmt.Sprintf("%s%s", fieldLine, crlf)), false, nil
}
