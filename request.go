package gomap

import "github.com/cwinters8/gomap/arguments"

type Request[A arguments.Args] struct {
	Using []Capability     `json:"using"`
	Calls []*Invocation[A] `json:"methodCalls"`
}
