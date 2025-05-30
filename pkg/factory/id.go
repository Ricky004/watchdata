package factory

import (
	"fmt"
	"regexp"
)

var (
	// idRegex is a regex that match the valid id.
	idRegex = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,30}$`)
)

type Id struct {
	id string
}

func (id Id) String() string {
	return id.id
}

func NewId(id string) (Id, error) {
	if !idRegex.MatchString(id) {
		return Id{}, fmt.Errorf("not a valid factory name %q", id)
	}
	return Id{id: id}, nil
}

func MustNewId(id string) Id {
	i, err := NewId(id)
	if err != nil {
		panic(err)
	}
	return i
}
