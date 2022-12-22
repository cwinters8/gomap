package gomap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Session    *Session
	HTTPClient *http.Client
	token      string
}

func NewClient(sessionURL string, bearerToken string) (*Client, error) {
	c := Client{
		token:      bearerToken,
		HTTPClient: http.DefaultClient,
	}
	// get session
	status, body, err := c.makeRequest(http.MethodGet, sessionURL, nil)
	if err != nil {
		return nil, fmt.Errorf("request status %d\nfailed to make session request: %w", status, err)
	}
	if body == nil {
		return nil, fmt.Errorf("nil session body cannot be used")
	}
	var sess Session
	if err := json.Unmarshal(body, &sess); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session body: %w", err)
	}
	c.Session = &sess
	return &c, nil
}

type Session struct {
	PrimaryAccounts *Accounts `json:"primaryAccounts"`
	APIURL          string    `json:"apiURL"`
}

type Accounts struct {
	Core       string `json:"urn:ietf:params:jmap:core"`
	Mail       string `json:"urn:ietf:params:jmap:mail"`
	Submission string `json:"urn:ietf:params:jmap:submission"`
}

func (c *Client) makeRequest(method string, url string, body []byte) (int, []byte, error) {
	var (
		req *http.Request
		err error
	)
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
		if err != nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("failed to create new request: %w", err)
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("failed to create new request: %w", err)
		}
	}
	req.Header = http.Header{
		"Authorization": []string{
			fmt.Sprintf("Bearer %s", c.token),
		},
		"Content-Type": []string{
			"application/json",
		},
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		status := http.StatusInternalServerError
		if resp != nil && resp.StatusCode > 0 {
			status = resp.StatusCode
		}
		return status, nil, fmt.Errorf("failed to make %s request to %s: %w", req.Method, req.URL, err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, body, fmt.Errorf("failed to parse response body: %w", err)
	}
	return http.StatusOK, respBody, nil
}
