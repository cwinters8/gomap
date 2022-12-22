package gomap

type Method interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// TODO: need to figure out what methods a call value needs
type CallValue interface {
	~string | Query
}

type Filter struct {
	Name string `json:"name"`
}

type Query struct {
	AccountID string `json:"accountId"`
	Filter    Filter `json:"filter"`
}

// these constants are mainly here for reference
const (
	MethodGetMailbox         = "Mailbox/get"
	MethodQueryMailbox       = "Mailbox/query"
	MethodSetEmail           = "Email/set"
	MethodSetEmailSubmission = "EmailSubmission/set"
)
