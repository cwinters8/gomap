package methods

import (
	"fmt"
)

type Method[A Args] struct {
	Prefix string
	Type   MethodType
	Args   A
}

func (m Method[A]) Name() string {
	return fmt.Sprintf("%s/%s", m.Prefix, m.Type)
}

type MethodType string

const (
	QueryMethod MethodType = "query"
	GetMethod   MethodType = "get"
	SetMethod   MethodType = "set"
)
