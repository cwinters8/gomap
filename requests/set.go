package requests

import (
	"fmt"

	"github.com/cwinters8/gomap/objects"

	"github.com/google/uuid"
)

type Set struct {
	ID   uuid.UUID
	Body *SetBody
}

func NewSet(acctID string, newObject objects.Object) (Set, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return NilSet, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	return Set{
		ID: id,
		Body: &SetBody{
			AccountID: acctID,
			Create:    newObject,
		},
	}, nil
}

func (s Set) GetID() uuid.UUID {
	return s.ID
}

func (s Set) Name() string {
	return "set"
}

func (s Set) Method() (string, error) {
	if s.Body.Create == nil {
		return "", fmt.Errorf("s.Create must not be nil")
	}
	return fmt.Sprintf("%s/set", s.Body.Create.Name()), nil
}

func (s Set) BodyMap() map[string]any {
	return map[string]any{
		"accountId": s.Body.AccountID,
		"create":    s.Body.Create,
	}
}

func (s Set) MarshalJSON() ([]byte, error) {
	return marshalJSON(Call(s))
}

type SetBody struct {
	AccountID string
	Create    objects.Object
}

// empty instance of Set struct
var NilSet = Set{}
