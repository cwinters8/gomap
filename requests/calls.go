package requests

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Call struct {
	ID        uuid.UUID
	AccountID string
	Method    string
	Arguments map[string]any
	OnSuccess func(map[string]any) error
	// OnError   func(error) error
}

func (c Call) MarshalJSON() ([]byte, error) {
	if c.Arguments == nil {
		return nil, fmt.Errorf("c.Arguments field must not be nil")
	}
	c.Arguments["accountId"] = c.AccountID
	slice := [3]any{c.Method, c.Arguments, c.ID}
	return json.Marshal(slice)
}
