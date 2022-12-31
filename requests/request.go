package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/results"
	"github.com/cwinters8/gomap/utils"
)

type Request struct {
	Using        []Capability `json:"using"`
	Calls        []Call       `json:"methodCalls"`
	SessionState string       `json:"sessionState"`
}

func NewRequest(calls []Call) *Request {
	return &Request{
		Using: []Capability{
			UsingCore,
			UsingMail,
		},
		Calls: calls,
	}
}

// adds the submission capability to the r.Using slice
func (r *Request) UseSubmission() {
	r.Using = append(r.Using, UsingSubmission)
}

func (r *Request) Send(c *client.Client) (*results.Results, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json from request: %w", err)
	}
	status, body, err := c.HttpRequest(http.MethodPost, c.Session.APIURL, b)
	if err != nil {
		return nil, fmt.Errorf("status %d - http request failed: %w", status, err)
	}
	if os.Getenv("RESPONSE_DEBUG") == "true" {
		// write raw body to file to allow for examination of full response
		if err := utils.WriteJSON("jmap_response", "../tmp/responses", body); err != nil {
			fmt.Printf("warning: failed to write json response to file: %s\n", err.Error())
		}
	}
	var results results.Results
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if len(results.Errors) > 0 {
		return nil, fmt.Errorf("found method errors: `%s`", utils.Prettier(results.Errors))
	}
	return &results, nil
}
