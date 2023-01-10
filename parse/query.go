package parse

import (
	"encoding/json"
	"fmt"
)

func QueryResponseBody(body map[string]any) (ids []string, err error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body to json: %w", err)
	}
	var resp queryResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query response: %w", err)
	}
	return resp.IDs, nil
}

type queryResponse struct {
	IDs []string `json:"ids"`
}
