package emails

import (
	"encoding/json"
	"fmt"

	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/results"
)

type Get struct {
	requests.Get
	FetchTextBody bool
	FetchHTMLBody bool
}

func (g Get) MarshalJSON() ([]byte, error) {
	m := g.BodyMap()
	m["fetchTextBodyValues"] = g.FetchTextBody
	m["fetchHTMLBodyValues"] = g.FetchHTMLBody
	return requests.MarshalCall(g.ID, "Email/get", m)
}

type Result struct {
	results.Get
	Body *ResultBody
}

type ResultBody struct {
	results.GetBody
	List []Email `json:"list"`
}

func (r *Result) Parse(rawBody any) error {
	b, err := results.ParseBytes(rawBody)
	if err != nil {
		return fmt.Errorf("failed to parse bytes from raw body: %w", err)
	}
	var body ResultBody
	if err := json.Unmarshal(b, &body); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}
	r.Body = &body
	return nil
}

func ParseGetResults(b []byte) ([]Email, error) {
	var result Result
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal email results: %w", err)
	}
	return result.Body.List, nil
}
