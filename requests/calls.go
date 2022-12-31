package requests

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Call interface {
	GetID() uuid.UUID
	Name() string
	Method() (string, error)
	BodyMap() map[string]any
	MarshalJSON() ([]byte, error)
}

type Filter struct {
	Name string `json:"name"`
}

func marshalJSON(c Call) ([]byte, error) {
	method, err := c.Method()
	if err != nil {
		return nil, fmt.Errorf("failed to get call method: %w", err)
	}
	s := [3]any{
		method,
		c.BodyMap(),
		c.GetID(),
	}
	return json.Marshal(s)
}
