package emails

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

func GetEmails(c *client.Client, emailIDs []string) (found []*Email, notFound []string, err error) {
	call, err := GetCall(c.Session.PrimaryAccounts.Mail, emailIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to construct Get call: %w", err)
	}
	responses, err := requests.Request(c, []*requests.Call{call}, false)
	if err != nil {
		return nil, nil, fmt.Errorf("request failure: %w", err)
	}
	return ParseRawResponseBody(responses[0].Body)
}

func GetCall(acctID string, emailIDs []string) (*requests.Call, error) {
	callID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	return &requests.Call{
		ID:        callID,
		AccountID: acctID,
		Method:    "Email/get",
		Arguments: map[string]any{
			"ids": emailIDs,
			"properties": []string{
				"mailboxIds",
				"from",
				"to",
				"subject",
				"bodyValues",
				"bodyStructure",
			},
			"bodyProperties":      []string{"type"},
			"fetchTextBodyValues": true,
			"fetchHTMLBodyValues": true,
		},
	}, nil
}

func ParseRawResponseBody(body map[string]any) (found []*Email, notFound []string, err error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal response body to json: %w", err)
	}
	var respBody responseBody
	if err := json.Unmarshal(b, &respBody); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	emails := []*Email{}
	for _, rawEmail := range respBody.List {
		emails = append(emails, &Email{
			ID:         rawEmail.ID,
			MailboxIDs: rawEmail.MailboxIDs.IDs,
			From:       rawEmail.From,
			To:         rawEmail.To,
			Subject:    rawEmail.Subject,
			Body: &Body{
				Value: rawEmail.BodyValue.Value,
				Type:  BodyType(rawEmail.BodyStructure.Type),
			},
		})
	}
	return emails, respBody.NotFound, nil
}

func (e *Email) Get(acctID string) (*requests.Call, error) {
	if len(e.ID) < 1 {
		return nil, fmt.Errorf("e.ID field must be populated")
	}
	call, err := GetCall(acctID, []string{e.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to construct Get call: %w", err)
	}
	call.OnSuccess = func(m map[string]any) error {
		found, _, err := ParseRawResponseBody(m)
		if err != nil {
			return fmt.Errorf("failed to parse raw response body: %w", err)
		}
		if len(found) < 1 {
			return fmt.Errorf("email id %s not found", e.ID)
		}
		email := found[0]
		e.MailboxIDs = email.MailboxIDs
		e.From = email.From
		e.To = email.To
		e.Subject = email.Subject
		e.Body = email.Body
		return nil
	}
	return call, nil
}

type responseBody struct {
	List     []*result `json:"list"`
	NotFound []string  `json:"notFound"`
}

type result struct {
	ID            string         `json:"id"`
	MailboxIDs    *mailboxes     `json:"mailboxIds"`
	From          []*Address     `json:"from"`
	To            []*Address     `json:"to"`
	Subject       string         `json:"subject"`
	BodyValue     *bodyValue     `json:"bodyValues"`
	BodyStructure *bodyStructure `json:"bodyStructure"`
}

type mailboxes struct {
	IDs []string
}

func (m *mailboxes) UnmarshalJSON(b []byte) error {
	var raw map[string]bool
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal mailbox ids to map: %w", err)
	}
	boxIDs := []string{}
	for k, v := range raw {
		if v {
			boxIDs = append(boxIDs, k)
		}
	}
	m.IDs = boxIDs
	return nil
}

type bodyValue struct {
	Value string
}

func (v *bodyValue) UnmarshalJSON(b []byte) error {
	var raw map[string]map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal body values to map: %w", err)
	}
	values := make([]string, len(raw))
	for k, v := range raw {
		key, err := strconv.Atoi(k)
		if err != nil {
			return fmt.Errorf("failed to convert key `%s` to int: %w", k, err)
		}
		val, ok := v["value"].(string)
		if !ok {
			return fmt.Errorf("failed to cast value to string. %s", utils.Describe(v["value"]))
		}
		values[key-1] = val
	}
	v.Value = strings.Join(values, " ")
	return nil
}

type bodyStructure struct {
	Type string `json:"type"`
}
