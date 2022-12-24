package gomap

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cwinters8/gomap/arguments"
	"github.com/cwinters8/gomap/utils"
)

type Request[A arguments.Args] struct {
	Using        []Capability     `json:"using"`
	Calls        []*Invocation[A] `json:"methodCalls"`
	SessionState string           `json:"sessionState"`
}

func NewRequest[A arguments.Args](calls []*Invocation[A]) *Request[A] {
	return &Request[A]{
		Using: []Capability{
			UsingCore,
			UsingMail,
		},
		Calls: calls,
	}
}

// adds the submission capability to the r.Using slice
func (r *Request[A]) UseSubmission() {
	r.Using = append(r.Using, UsingSubmission)
}

func (r *Request[A]) Send(c *Client) (*Response[A], error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json from request: %w", err)
	}
	status, body, err := c.httpRequest(http.MethodPost, c.Session.APIURL, b)
	if err != nil {
		return nil, fmt.Errorf("status %d - http request failed: %w", status, err)
	}
	var resp Response[A]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	var errs []Error
	// check for method errors
	for _, i := range resp.Results {
		if i.Method.Err != nil {
			errs = append(errs, *i.Method.Err)
		}
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("found method errors: %s", utils.Prettier(errs))
	}
	return &resp, nil
}

type Response[A arguments.Args] struct {
	Results []*Invocation[A] `json:"methodResponses"`
}
