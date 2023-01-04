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
	OnSuccess func([]byte) error
	OnError   func(error) error
}

func (c Call) MarshalJSON() ([]byte, error) {
	c.Arguments["accountId"] = c.AccountID
	slice := [3]any{c.Method, c.Arguments, c.ID}
	return json.Marshal(slice)
}

func MarshalCall(id uuid.UUID, method string, body map[string]any) ([]byte, error) {
	s := [3]any{method, body, id}
	return json.Marshal(s)
}
