package results

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Get struct {
	ID     uuid.UUID
	Prefix string
	Body   *GetBody
}

func (g *Get) Name() string {
	return "get"
}
func (g *Get) Method() (string, error) {
	if len(g.Prefix) < 1 {
		return "", fmt.Errorf("g.Prefix must not be empty")
	}
	return fmt.Sprintf("%s/get", g.Prefix), nil
}
func (g *Get) Parse(rawBody any) error {
	b, err := parseBytes(rawBody)
	if err != nil {
		return fmt.Errorf("failed to parse bytes from raw body: %w", err)
	}
	var body GetBody
	if err := json.Unmarshal(b, &body); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}
	g.Body = &body
	return nil
}

type GetBody struct {
	AccountID string `json:"accountId"`
	// List      []objects.Object `json:"list"`
	NotFound []string `json:"notFound"`
}
