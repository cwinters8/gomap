package results

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Set struct {
	Created   *Created `json:"created"`
	Updated   *Updated `json:"updated"`
	AccountID string   `json:"accountId"`
	OldState  string   `json:"oldState"`
	NewState  string   `json:"newState"`
}

type Created struct {
	ID uuid.UUID `json:"-"`
	created
}

func (c *Created) UnmarshalJSON(b []byte) error {
	var m map[string]created
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unmarshal json to map: %w", err)
	}
	for k, v := range m {
		id, err := uuid.Parse(k)
		if err != nil {
			return fmt.Errorf("failed to parse key `%s` as uuid: %w", k, err)
		}
		c.ID = id
		c.created = v
		break // map should have a single key, so only need first iteration
	}
	return nil
}

type created struct {
	ServerID string `json:"id"`
	BlobID   string `json:"blobId"`
	ThreadID string `json:"threadId"`
	Size     int    `json:"size"`
}

type Updated struct {
	ServerID string `json:"-"`
}

func (u *Updated) UnmarshalJSON(b []byte) error {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unmarshal json to map: %w", err)
	}
	for k := range m {
		u.ServerID = k
		break // map should have a single key, so only need first iteration
	}
	return nil
}
