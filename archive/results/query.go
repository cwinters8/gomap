package results

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Query struct {
	ID     uuid.UUID
	Prefix string
	Body   *QueryBody
}

func (q *Query) Name() string {
	return "query"
}

func (q *Query) Method() (string, error) {
	if len(q.Prefix) < 1 {
		return "", fmt.Errorf("q.Prefix must not be empty")
	}
	return fmt.Sprintf("%s/query", q.Prefix), nil
}

func (q *Query) Parse(rawBody any) error {
	b, err := parseBytes(rawBody)
	if err != nil {
		return fmt.Errorf("failed to parse bytes from raw body: %w", err)
	}
	var body QueryBody
	if err := json.Unmarshal(b, &body); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}
	q.Body = &body
	return nil
}

type QueryBody struct {
	AccountID string   `json:"accountId"`
	IDs       []string `json:"ids"`
	Total     int      `json:"total"`
	Filter    *Filter  `json:"filter"`
}

type Filter struct {
	Name string `json:"name"`
}
