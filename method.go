package gomap

import (
	"encoding/json"
	"fmt"
)

// implementations of Arguments should include appropriate json tags
type Arguments interface {
	*MailboxQuery
	GetMethod() Method
}

type Invocation[A Arguments] struct {
	ID     string
	Method Method
	Args   A
}

type MailboxQuery struct {
	AccountID string `json:"accountId"`
	Filter    Filter `json:"filter"`
}

func (mq *MailboxQuery) GetMethod() Method {
	return "Mailbox/query"
}

type Filter struct {
	Name string `json:"name"`
}

type Method string

// these constants are mainly here for reference
const (
	GetMailbox         Method = "Mailbox/get"
	SetEmail           Method = "Email/set"
	SetEmailSubmission Method = "EmailSubmission/set"
)

func (i Invocation[A]) MarshalJSON() ([]byte, error) {
	raw := []any{
		i.Method,
		i.Args,
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
	i.Method = Method(method)
	args, err := json.Marshal(raw[1])
	if err != nil {
		return fmt.Errorf("failed to marshal raw args back to json: %w", err)
	}
	var a A
	if err := json.Unmarshal(args, &a); err != nil {
		return fmt.Errorf("failed to unmarshal args: %w", err)
	}
	i.Args = a
	id, ok := raw[2].(string)
	if !ok {
		return fmt.Errorf("failed to coerce id to string")
	}
	i.ID = id
	return nil
}
