package requests

import (
	"fmt"

	"github.com/google/uuid"
)

type Query struct {
	ID     uuid.UUID
	Prefix string
	Body   *QueryBody
}

// creates a new Query that will filter objects by name
func NewQuery(acctID string, prefix string, name string) (Query, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return NilQuery, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	return Query{
		ID:     id,
		Prefix: prefix,
		Body: &QueryBody{
			AccountID: acctID,
			Filter: &Filter{
				Name: name,
			},
		},
	}, nil
}

func (q Query) GetID() uuid.UUID {
	return q.ID
}

func (q Query) Name() string {
	return "query"
}

func (q Query) Method() (string, error) {
	if len(q.Prefix) < 1 {
		return "", fmt.Errorf("q.Prefix must not be empty")
	}
	return fmt.Sprintf("%s/query", q.Prefix), nil
}

func (q Query) BodyMap() map[string]any {
	return map[string]any{
		"accountId": q.Body.AccountID,
		"filter":    q.Body.Filter,
	}
}

func (q Query) MarshalJSON() ([]byte, error) {
	return marshalJSON(Call(q))
}

type QueryBody struct {
	AccountID string
	Filter    *Filter
}

// empty instance of Query struct
var NilQuery = Query{}
