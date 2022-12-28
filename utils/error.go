package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var ErrNotImplemented = fmt.Errorf("not implemented")

// returns a string with the value and its underlying type
//
// useful for providing meaningful error messages when a type coercion fails
func Describe(value any) string {
	if value == nil {
		return "value is nil"
	}
	return fmt.Sprintf("value `%v` has underlying type `%s`", Prettier(value), reflect.TypeOf(value).Kind().String())
}

// provides a more human readable string representation of v
func Prettier(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("warning: failed to marshal v to formatted output: %s\n", err.Error())
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
