package utils

import (
	"fmt"
	"testing"
)

// logs a formatted string, then fails the test immediately
func Failf(t *testing.T, format string, args ...any) {
	t.Logf(fmt.Sprintf("%s\n", format), args...)
	t.FailNow()
}

// if b is true, causes the test to fail with the supplied msg
func Checkf(t *testing.T, b bool, format string, args ...any) {
	if b {
		t.Errorf(format, args...)
	}
}

type Case struct {
	Check  bool // case fails if check is true
	Format string
	Args   []any
}
