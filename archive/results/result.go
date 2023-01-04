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
}

type Error struct {
	ID          uuid.UUID      `json:"-"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
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
		case "description":
			d, ok := v.(string)
			if !ok {
				return fmt.Errorf("failed to cast description to string. %s", utils.Describe(v))
			}
			e.Description = d
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
			var result Result
			switch action {
			case "query":
				result = &Query{
					ID:     id,
					Prefix: methodTypes[0],
				}
			case "set":
				result = &Set{
					ID:     id,
					Prefix: methodTypes[0],
				}
			case "get":
				result = &Get{
					ID:     id,
					Prefix: methodTypes[0],
				}
			default:
				return fmt.Errorf("unsupported action `%s`", action)
			}
			if err := result.Parse(resp[1]); err != nil {
				return fmt.Errorf("failed to parse body: %w", err)
			}
			r.Results = append(r.Results, result)
		}
	}
	return nil
}

func parseBytes(raw any) ([]byte, error) {
	var b []byte
	if str, ok := raw.(string); !ok {
		jsonBytes, err := json.Marshal(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal raw body to json: %w", err)
		}
		b = jsonBytes
	} else {
		b = []byte(str)
	}
	return b, nil
}
func ParseBytes(raw any) ([]byte, error) {
	var b []byte
	if str, ok := raw.(string); !ok {
		jsonBytes, err := json.Marshal(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal raw body to json: %w", err)
		}
		b = jsonBytes
	} else {
		b = []byte(str)
	}
	return b, nil
}
