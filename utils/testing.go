package utils

import (
	"fmt"
)

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

type Cases []*Case

// runs onCheck when a case's Check is true
func (cases Cases) Iterator(onCheck func(c *Case)) {
	passed := 0
	failed := 0
	for _, c := range cases {
		if c.Check {
			onCheck(c)
			failed++
			continue
		}
		passed++
	}
	fmt.Printf("processed %d total cases: %d passed; %d failed\n", passed+failed, passed, failed)
}

func (cases *Cases) Append(c ...*Case) {
	caseSlice := []*Case(*cases)
	*cases = append(caseSlice, c...)
}
