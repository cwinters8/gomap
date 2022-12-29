package requests

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cwinters8/gomap/requests/arguments"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

type Invocation[A arguments.Args] struct {
	ID     uuid.UUID
	Method *Method[A]
}

type Method[A arguments.Args] struct {
	Prefix string
	Type   MethodType
	Args   A
	Err    *Error
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

type Error struct {
	Type        string         `json:"type"`
	Properties  []string       `json:"properties"`
	MoreDetails map[string]any `json:"moreDetails"`
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
	m := Method[A]{
		Prefix: strings.Split(method, "/")[0],
	}
	if method == "error" {
		var methodErr Error
		if err := json.Unmarshal(args, &methodErr); err != nil {
			return fmt.Errorf("failed to unmarshal error: %w", err)
		}
		m.Err = &methodErr
	} else {
		var a A
		if err := json.Unmarshal(args, &a); err != nil {
			return fmt.Errorf("failed to unmarshal args: %w", err)
		}
		m.Args = a
	}
	i.Method = &m
	id, ok := raw[2].(string)
	if !ok {
		return fmt.Errorf("failed to coerce id to string")
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse uuid: %w", err)
	}
	i.ID = parsedID
	return nil
}

func (e *Error) UnmarshalJSON(b []byte) error {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unmarshal error: %w", err)
	}
	details := map[string]any{}
	for k, v := range m {
		switch k {
		case "type":
			t, ok := v.(string)
			if !ok {
				return fmt.Errorf("failed to coerce type to string. %s", utils.Describe(v))
			}
			e.Type = t
		case "properties":
			p, ok := v.([]string)
			if !ok {
				return fmt.Errorf("failed to coerce properties to []string. %s", utils.Describe(v))
			}
			e.Properties = p
		default:
			details[k] = v
		}
	}
	e.MoreDetails = details
	return nil
}
