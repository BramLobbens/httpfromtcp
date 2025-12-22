package headers

import "fmt"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	err = fmt.Errorf("error: not yet implemented")
	return 0, false, err
}
