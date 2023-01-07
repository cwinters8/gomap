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

func GetEmails(c *client.Client, emailIDs []string) ([]*Email, error) {
	callID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	call := requests.Call{
		ID:        callID,
		AccountID: c.Session.PrimaryAccounts.Mail,
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
	}
	responses, err := requests.Request(c, []*requests.Call{&call}, false)
	if err != nil {
		return nil, fmt.Errorf("request failure: %w", err)
	}
	resp := responses[0]
	b, err := json.Marshal(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response body to json: %w", err)
	}
	var body responseBody
	if err := json.Unmarshal(b, &body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	emails := []*Email{}
	for _, rawEmail := range body.List {
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
	if len(body.NotFound) > 0 {
		fmt.Printf("warning: some email IDs were not found: %v\n", body.NotFound)
	}
	return emails, nil
}

func (e *Email) Get(acctID string) (*requests.Call, error) {
	if len(e.ID) < 1 {
		return nil, fmt.Errorf("e.ID field must be populated")
	}
	callID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	return &requests.Call{
		ID:        callID,
		AccountID: acctID,
		Method:    "Email/get",
		Arguments: map[string]any{
			"ids": []string{e.ID},
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
		OnSuccess: func(m map[string]any) error {
			list, ok := m["list"].([]any)
			if !ok {
				return fmt.Errorf("failed to cast result list to slice. %s", utils.Describe(m["list"]))
			}
			result, ok := list[0].(map[string]any)
			if !ok {
				return fmt.Errorf("failed to cast result value to map. %s", utils.Describe(list[0]))
			}
			id, ok := result["id"].(string)
			if !ok {
				return fmt.Errorf("failed to cast id to string. %s", utils.Describe(m["id"]))
			}
			if id != e.ID {
				return fmt.Errorf("returned id does not match requested email id. wanted %s; got %s", e.ID, id)
			}
			rawBoxIDs, ok := result["mailboxIds"].(map[string]any)
			if !ok {
				return fmt.Errorf("failed to coerce mailbox IDs to map. %s", utils.Describe(result["mailboxIds"]))
			}
			boxIDs := []string{}
			for k := range rawBoxIDs {
				boxIDs = append(boxIDs, k)
			}
			e.MailboxIDs = boxIDs
			rawFrom, ok := result["from"].([]any)
			if !ok {
				return fmt.Errorf("failed to cast from addresses to slice. %s", utils.Describe(result["from"]))
			}
			from, err := parseAddresses(rawFrom)
			if err != nil {
				return fmt.Errorf("failed to parse from addresses: %w", err)
			}
			e.From = from
			rawTo, ok := result["to"].([]any)
			if !ok {
				return fmt.Errorf("failed to cast to addresses to slice. %s", utils.Describe(result["to"]))
			}
			to, err := parseAddresses(rawTo)
			if err != nil {
				return fmt.Errorf("failed to parse to addresses: %w", err)
			}
			e.To = to
			subj, ok := result["subject"].(string)
			if !ok {
				return fmt.Errorf("failed to cast subject to string. %s", utils.Describe(result["subject"]))
			}
			e.Subject = subj
			rawBodyValues, ok := result["bodyValues"].(map[string]any)
			if !ok {
				return fmt.Errorf("failed to cast body values to map. %s", utils.Describe(result["bodyValues"]))
			}
			values := make([]string, len(rawBodyValues))
			for k, v := range rawBodyValues {
				key, err := strconv.Atoi(k)
				if err != nil {
					return fmt.Errorf("failed to convert key `%s` to int: %w", k, err)
				}
				rawValue, ok := v.(map[string]any)
				if !ok {
					return fmt.Errorf("failed to cast raw value to map. %s", utils.Describe(v))
				}
				val, ok := rawValue["value"].(string)
				if !ok {
					return fmt.Errorf("failed to cast value to string. %s", utils.Describe(rawValue["value"]))
				}
				values[key-1] = val
			}
			structure, ok := result["bodyStructure"].(map[string]any)
			if !ok {
				return fmt.Errorf("failed to cast body structure to map. %s", utils.Describe(result["bodyStructure"]))
			}
			t, ok := structure["type"].(string)
			if !ok {
				return fmt.Errorf("failed to cast body type to string. %s", utils.Describe(structure["type"]))
			}
			e.Body = &Body{
				Value: strings.Join(values, " "),
				Type:  BodyType(t),
			}
			return nil
		},
	}, nil
}

func parseAddresses(raw []any) ([]*Address, error) {
	addresses := []*Address{}
	for _, val := range raw {
		address, ok := val.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed to cast address to map. %s", utils.Describe(val))
		}
		email, ok := address["email"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to cast email to string. %s", utils.Describe(address["email"]))
		}
		name, ok := address["name"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to cast name to string. %s", utils.Describe(address["name"]))
		}
		addresses = append(addresses, &Address{
			Email: email,
			Name:  name,
		})
	}
	return addresses, nil
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
