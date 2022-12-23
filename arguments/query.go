package arguments

type Query struct {
	AccountID string `json:"accountId"`
	Filter    Filter `json:"filter"`
}

type QueryResp struct {
	IDs []string `json:"ids"`
}
