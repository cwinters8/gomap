package results

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

type Result interface {
	Name() string
	Method() (string, error)
	Parse(any) error
	UnmarshalJSON([]byte) error
}

type Error struct {
	ID          uuid.UUID      `json:"-"`
	Type        string         `json:"type"`
	Properties  []string       `json:"properties"`
	MoreDetails map[string]any `json:"moreDetails"`
}

func (e *Error) Parse(a any) error {
	m, ok := a.(map[string]any)
	if !ok {
		return fmt.Errorf("failed to coerce error body to map. %s", utils.Describe(a))
	}
	for k, v := range m {
		switch k {
		case "type":
			t, ok := v.(string)
			if !ok {
				return fmt.Errorf("failed to coerce type to string. %s", utils.Describe(m["type"]))
			}
			e.Type = t
		case "properties":
			props, ok := m["properties"].([]any)
			if !ok {
				return fmt.Errorf("failed to coerce properties to slice. %s", utils.Describe(m["properties"]))
			}
			for _, p := range props {
				prop, ok := p.(string)
				if !ok {
					return fmt.Errorf("failed to coerce property to slice. %s", utils.Describe(p))
				}
				e.Properties = append(e.Properties, prop)
			}
		}
	}
	return nil
}

type Results struct {
	Results      []Result
	Errors       []*Error
	SessionState string
}

type rawResults struct {
	Responses    [][3]any `json:"methodResponses"`
	SessionState string   `json:"sessionState"`
}

func (r *Results) UnmarshalJSON(b []byte) error {
	var raw rawResults
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal raw results: %w", err)
	}
	r.SessionState = raw.SessionState
	for _, resp := range raw.Responses {
		method, ok := resp[0].(string)
		if !ok {
			return fmt.Errorf("failed to coerce method to string. %s", utils.Describe(resp[0]))
		}
		rawID, ok := resp[2].(string)
		if !ok {
			return fmt.Errorf("failed to coerce id to string. %s", utils.Describe(resp[2]))
		}
		id, err := uuid.Parse(rawID)
		if err != nil {
			return fmt.Errorf("failed to parse uuid: %w", err)
		}
		methodTypes := strings.Split(method, "/")
		if methodTypes[0] == "error" {
			// handle error result
			e := Error{
				ID: id,
			}
			if err := e.Parse(resp[1]); err != nil {
				return fmt.Errorf("failed to parse error: %w", err)
			}
			r.Errors = append(r.Errors, &e)
		} else {
			if len(methodTypes) < 2 {
				return fmt.Errorf("unsupported method `%s`. non-error methods must be in the format `Type/action`", method)
			}
			action := methodTypes[1]

			switch action {
			// TODO: handle query actions
			case "set":
				s := Set{
					ID:     id,
					Prefix: methodTypes[0],
				}
				if err := s.Parse(resp[1]); err != nil {
					return fmt.Errorf("failed to parse set body: %w", err)
				}
				r.Results = append(r.Results, &s)
			default:
				return fmt.Errorf("unsupported action `%s`", action)
			}
		}
	}
	return nil
}
