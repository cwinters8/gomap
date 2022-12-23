package gomap

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cwinters8/gomap/arguments"
)

type Invocation[A arguments.Args] struct {
	ID     string
	Method *Method[A]
}

type Method[A arguments.Args] struct {
	Prefix string
	Args   A
}

func (m Method[A]) Name() string {
	return fmt.Sprintf("%s/query", m.Prefix)
}

func (i Invocation[A]) MarshalJSON() ([]byte, error) {
	raw := []any{
		// method
		i.Method.Name(),

		// args
		i.Method.Args,

		// id
		i.ID,
	}
	return json.Marshal(raw)
}

func (i *Invocation[A]) UnmarshalJSON(b []byte) error {
	var raw []any
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal raw json: %w", err)
	}
	method, ok := raw[0].(string)
	if !ok {
		return fmt.Errorf("failed to coerce method to string")
	}
	args, err := json.Marshal(raw[1])
	if err != nil {
		return fmt.Errorf("failed to marshal raw args back to json: %w", err)
	}
	var a A
	if err := json.Unmarshal(args, &a); err != nil {
		return fmt.Errorf("failed to unmarshal args: %w", err)
	}
	m := Method[A]{
		Prefix: strings.Split(method, "/")[0],
		Args:   a,
	}
	i.Method = &m
	id, ok := raw[2].(string)
	if !ok {
		return fmt.Errorf("failed to coerce id to string")
	}
	i.ID = id
	return nil
}
