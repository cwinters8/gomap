package utils

import "fmt"

type Case struct {
	Check   bool // case fails if check is true
	Message string
	Args    []any // deprecated
}

func NewCase(check bool, msg string, args ...any) *Case {
	return &Case{
		Check:   check,
		Message: fmt.Sprintf(msg, args...),
	}
}
