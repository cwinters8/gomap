package requests

import (
	"fmt"

	"github.com/google/uuid"
)

type Get struct {
	ID     uuid.UUID
	Prefix string
	Body   *GetBody
}

func NewGet(acctID, prefix string, ids, properties []string) (Get, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return NilGet, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	body := GetBody{AccountID: acctID}
	if len(ids) > 0 {
		body.IDs = ids
	}
	if len(properties) > 0 {
		body.Properties = properties
	}
	return Get{
		ID:     id,
		Prefix: prefix,
		Body:   &body,
	}, nil
}

func (g Get) GetID() uuid.UUID {
	return g.ID
}

func (g Get) Name() string {
	return "get"
}

func (g Get) Method() (string, error) {
	if len(g.Prefix) < 1 {
		return "", fmt.Errorf("g.Prefix must not be empty")
	}
	return fmt.Sprintf("%s/get", g.Prefix), nil
}

func (g Get) BodyMap() map[string]any {
	m := map[string]any{
		"accountId": g.Body.AccountID,
	}
	if g.Body.IDs != nil && len(g.Body.IDs) > 0 {
		m["ids"] = g.Body.IDs
	}
	if g.Body.Properties != nil && len(g.Body.Properties) > 0 {
		m["properties"] = g.Body.Properties
	}
	return m
}

func (g Get) MarshalJSON() ([]byte, error) {
	return marshalJSON(g)
}

type GetBody struct {
	AccountID  string
	IDs        []string
	Properties []string
}

var NilGet = Get{}
