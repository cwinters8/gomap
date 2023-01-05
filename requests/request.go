package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/utils"
)

func Request(c *client.Client, calls []*Call, usingSubmission bool) error {
	using := []Capability{UsingCore, UsingMail}
	if usingSubmission {
		using = append(using, UsingSubmission)
	}
	r := Req{
		Using: using,
		Calls: calls,
	}
	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("failed to marshal request to json: %w", err)
	}
	status, result, err := c.HttpRequest(http.MethodPost, c.Session.APIURL, b)
	if err != nil {
		return fmt.Errorf("status %d - http request failure: %w", status, err)
	}
	if os.Getenv("REQUEST_DEBUG") == "true" {
		// write raw json to file to allow for examination of full request and response
		raw := map[string][]byte{
			"request":  b,
			"response": result,
		}
		if err := utils.WriteJSON("jmap", "../tmp/raw", raw); err != nil {
			fmt.Printf("warning: failed to write json response to file: %s\n", err.Error())
		}
	}
	var resp Resp
	if err := json.Unmarshal(result, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	errs := []map[string]any{}
Responses:
	for _, r := range resp.MethodResponses {
		method, ok := r[0].(string)
		if !ok {
			return fmt.Errorf("failed to cast method as string. %s", utils.Describe(r[0]))
		}
		body, ok := r[1].(map[string]any)
		if !ok {
			return fmt.Errorf("failed to cast response body to map. %s", utils.Describe(r[1]))
		}
		idStr, ok := r[2].(string)
		if !ok {
			return fmt.Errorf("failed to cast id to string. %s", utils.Describe(r[2]))
		}
		if method == "error" {
			body["id"] = idStr
			errs = append(errs, body)
			continue
		}
		for _, c := range calls {
			if c.ID.String() == idStr && c.Method == method {
				if err := c.OnSuccess(body); err != nil {
					return fmt.Errorf("call to OnSuccess failed: %w", err)
				}
				continue Responses
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("found method errors: %v", errs)
	}
	return nil
}

type Req struct {
	Using []Capability `json:"using"`
	Calls []*Call      `json:"methodCalls"`
}

type Resp struct {
	MethodResponses [][3]any `json:"methodResponses"`
}
