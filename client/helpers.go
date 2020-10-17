package client

import (
	"bytes"
	"strconv"
	"strings"
)

func joinIntSlice(s []int) string {
	out := new(strings.Builder)
	for i, n := range s {
		if i != len(s)-1 {
			out.WriteString(strconv.Itoa(n) + ",")
		} else {
			out.WriteString(strconv.Itoa(n))
		}
	}

	return out.String()
}

func fixBadArraySerialization(body []byte, fields [][]byte) []byte {
	var old, new []byte

	for _, field := range fields {
		old = append(append([]byte(`"`), field...), []byte(`":{}`)...)
		new = append(append([]byte(`"`), field...), []byte(`":[]`)...)
		body = bytes.ReplaceAll(body, old, new)
	}

	return body
}
