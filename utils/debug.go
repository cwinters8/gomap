package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func WriteJSON(name, dir string, raw []byte) error {
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		return fmt.Errorf("failed to unmarshal raw body: %w", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create dir `%s`: %w", dir, err)
	}
	t := time.Now()
	d := document{
		Timestamp: t,
		Body:      body,
	}
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal document to json: %w", err)
	}
	// timestamp format is RFC3339 with more precision
	filename := fmt.Sprintf("%s/%s-%s.json", dir, name, t.Format("2006-01-02T15:04:05.99Z07:00"))
	if err := os.WriteFile(filename, b, 0644); err != nil {
		return fmt.Errorf("failed to write to file `%s`: %w", filename, err)
	}
	return nil
}

type document struct {
	Timestamp time.Time      `json:"timestamp"`
	Body      map[string]any `json:"body"`
}
