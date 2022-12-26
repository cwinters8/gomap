package utils

import (
	"fmt"

	"testing"

	"github.com/joho/godotenv"
)

// logs a formatted string, then fails the test immediately
//
// TODO: deprecate this and replace with t.Fatalf
func Failf(t *testing.T, format string, args ...any) {
	t.Logf(fmt.Sprintf("%s\n", format), args...)
	t.FailNow()
}

// if b is true, causes the test to fail with the supplied msg
//
// TODO: deprecate this and iterate over slices of Case to reduce duplication
func Checkf(t *testing.T, b bool, format string, args ...any) {
	if b {
		t.Errorf(format, args...)
	}
}

// load env variables from .env
func Env(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Fatalf("failed to load .env: %s", err.Error())
	}
}

type Case struct {
	Check  bool // case fails if check is true
	Format string
	Args   []any
}
