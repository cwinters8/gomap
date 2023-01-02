package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func WriteJSON(name, dir string, raw map[string][]byte) error {
	result := map[string]any{}
	for k, v := range raw {
		var val any
		if err := json.Unmarshal(v, &val); err != nil {
			return fmt.Errorf("failed to unmarshal value for key `%s`: %w", k, err)
		}
		result[k] = val
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create dir `%s`: %w", dir, err)
	}
	t := time.Now()
	result["timestamp"] = t
	b, err := json.Marshal(result)
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
