package methods

type Query struct {
	AccountID string   `json:"accountId"`
	Filter    Filter   `json:"filter"`
	IDs       []string `json:"ids,omitempty"`
}
