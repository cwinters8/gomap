package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
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
	prefix := fmt.Sprintf("%s/%s-%s", dir, name, t.Format("2006-01-02T15:04:05.99Z07:00"))
	suffix := "json"
	filename := fmt.Sprintf("%s.%s", prefix, suffix)
	// check if file exists before writing
	for {
		if _, err := os.Stat(filename); err == nil {
			idx := len(prefix) - 1
			endStr := string(prefix[idx])
			end, err := strconv.Atoi(endStr)
			if err != nil {
				return fmt.Errorf("failed to convert character `%s` to int: %w", endStr, err)
			}
			switch end {
			case 0:
				prefix += "-1"
			default:
				prefix = fmt.Sprintf("%s%d", prefix[:idx-1], end+1)
			}
			filename = fmt.Sprintf("%s.%s", prefix, suffix)
		} else if errors.Is(err, os.ErrNotExist) {
			break
		} else {
			return fmt.Errorf("something went wrong with checking os.Stat on file `%s`: %w", filename, err)
		}
	}
	if err := os.WriteFile(filename, b, 0644); err != nil {
		return fmt.Errorf("failed to write to file `%s`: %w", filename, err)
	}
	return nil
}
