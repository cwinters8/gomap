package results

import (
	"encoding/json"
	"fmt"

	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

type Set struct {
	ID     uuid.UUID
	Prefix string
	Body   *SetBody
}

func (s *Set) Name() string {
	return "set"
}

func (s *Set) Method() (string, error) {
	if len(s.Prefix) < 1 {
		return "", fmt.Errorf("s.Prefix must not be empty")
	}
	return fmt.Sprintf("%s/set", s.Prefix), nil
}

func (s *Set) UnmarshalJSON(b []byte) error {
	return utils.ErrNotImplemented
}
func (s *Set) Parse(rawBody any) error {
	b, err := json.Marshal(rawBody)
	if err != nil {
		return fmt.Errorf("failed to marshal set body to json: %w", err)
	}
	var body SetBody
	if err := json.Unmarshal(b, &body); err != nil {
		return fmt.Errorf("failed  to unmarshal body: %w", err)
	}
	s.Body = &body
	return nil
}

type SetBody struct {
	Created    *Created    `json:"created"`
	NotCreated *NotCreated `json:"notCreated"`
	Updated    *Updated    `json:"updated"`
	AccountID  string      `json:"accountId"`
	OldState   string      `json:"oldState"`
	NewState   string      `json:"newState"`
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

type NotCreated struct {
	ID uuid.UUID `json:"-"`
	notCreated
}

func (nc *NotCreated) UnmarshalJSON(b []byte) error {
	var m map[string]notCreated
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unmarshal json to map: %w", err)
	}
	for k, v := range m {
		id, err := uuid.Parse(k)
		if err != nil {
			return fmt.Errorf("failed to parse key `%s` as uuid: %w", k, err)
		}
		nc.ID = id
		nc.notCreated = v
		break // map should have a single key, so only need first iteration
	}
	return nil
}

type notCreated struct {
	Properties []string `json:"properties"`
	Type       string   `json:"type"`
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
