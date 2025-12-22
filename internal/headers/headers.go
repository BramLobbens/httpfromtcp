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

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlf := []byte("\r\n")

	last := bytes.LastIndex(data, crlf)
	if last == -1 {
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
			return 0, false, fmt.Errorf("malformed header: %s", string(line))
		}

		if s := strings.TrimRight(string(parts[0]), " "); len(s) < len(parts[0]) {
			return 0, false, fmt.Errorf("invalid formatting in field-name")
		}

		key := string(bytes.TrimSpace(parts[0]))
		val := string(bytes.TrimSpace(parts[1]))

		if _, ok := h[key]; ok {
			return 0, false, fmt.Errorf("key-value pair already exists")
		}
		h[key] = val
		numBytesParsed += len(line) + len(crlf)
	}

	return numBytesParsed, false, nil
}
