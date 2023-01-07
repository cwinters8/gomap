package requests

import (
	"encoding/json"

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
	c.Arguments["accountId"] = c.AccountID
	slice := [3]any{c.Method, c.Arguments, c.ID}
	return json.Marshal(slice)
}
