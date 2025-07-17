package config

import (
	"fmt"
	"regexp"
)

var (
	UriRegex = regexp.MustCompile(`(?s:^(?P<Scheme>[A-Za-z][A-Za-z0-9+.-]+):(?P<Value>.*)$)`)
)

type Uri struct {
	scheme string
	value  string
}

func NewUri(input string) (Uri, error) {
	smatch := UriRegex.FindStringSubmatch(input)

	if len(smatch) != 3 {
		return Uri{}, fmt.Errorf("invalid uri: %q", input)
	}
	return Uri{
		scheme: smatch[1],
		value:  smatch[2],
	}, nil
}

func MustNewUri(input string) Uri {
	uri, err := NewUri(input)
	if err != nil {
		panic(err)
	}

	return uri
}

func (uri Uri) Scheme() string {
	return uri.scheme
}

func (uri Uri) Value() string {
	return uri.value
}
