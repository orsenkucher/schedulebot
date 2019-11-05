package creds

import (
	"fmt"
	"strings"
)

// CreateFrom encodes your token
func CreateFrom(token string) string {
	bytes := []byte(token)
	magic(bytes)
	encoded := fmt.Sprint(bytes)
	encoded = makeReplacements(encoded)
	return encoded
}

func makeReplacements(s string) string {
	for _, r := range []struct {
		from string
		to   string
	}{
		{from: "[", to: "Credential{"},
		{from: " ", to: ", "},
		{from: "]", to: "}"},
	} {
		s = strings.ReplaceAll(s, r.from, r.to)
	}
	return s
}
