package emails

import (
	"encoding/json"
	"fmt"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/requests"

	"github.com/google/uuid"
)

func (e *Email) Submit(c *client.Client, draftMailboxID, sentMailboxID string) (submissionID string, err error) {
	call, err := SubmitCall(e.RequestID, c.Session.PrimaryAccounts.Mail, e.ID, draftMailboxID, sentMailboxID)
	if err != nil {
		return "", fmt.Errorf("failed to construct Submit call: %w", err)
	}
	responses, err := requests.Request(c, []*requests.Call{call}, true)
	if err != nil {
		return "", fmt.Errorf("submit request failure: %w", err)
	}
	if len(responses) < 1 {
		return "", fmt.Errorf("no responses returned")
	}
	created, err := ParseSubmitResponseBody(e.RequestID, responses[0].Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}
	e.Keywords.Draft = false
	sentBoxFound := false
	if len(draftMailboxID) > 0 {
		for idx, box := range e.MailboxIDs {
			switch box {
			case draftMailboxID:
				e.MailboxIDs = append(e.MailboxIDs[0:idx], e.MailboxIDs[idx+1:]...)
			case sentMailboxID:
				sentBoxFound = true
			}
		}
	}
	if len(sentMailboxID) > 0 && !sentBoxFound {
		e.MailboxIDs = append(e.MailboxIDs, sentMailboxID)
	}
	return created, nil
}

func SubmitCall(requestID uuid.UUID, acctID, emailID, draftMailboxID, sentMailboxID string) (*requests.Call, error) {
	callID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	args := map[string]any{
		"create": map[string]map[string]string{
			requestID.String(): {
				"emailId": emailID,
			},
		},
	}

	onSuccess := map[string]any{}
	if len(draftMailboxID) > 0 {
		onSuccess["keywords/$draft"] = nil
		onSuccess[fmt.Sprintf("mailboxIds/%s", draftMailboxID)] = nil
	}
	if len(sentMailboxID) > 0 {
		onSuccess[fmt.Sprintf("mailboxIds/%s", sentMailboxID)] = true
	}
	if len(onSuccess) > 0 {
		key := fmt.Sprintf("#%s", requestID.String())
		args["onSuccessUpdateEmail"] = map[string]map[string]any{key: onSuccess}
	}
	return &requests.Call{
		ID:        callID,
		AccountID: acctID,
		Method:    "EmailSubmission/set",
		Arguments: args,
	}, nil
}

func ParseSubmitResponseBody(requestID uuid.UUID, body map[string]any) (createdID string, err error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body to json: %w", err)
	}
	var resp setResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal submission response: %w", err)
	}
	if len(resp.NotCreated) > 0 {
		if failure, ok := resp.NotCreated[requestID.String()]; ok {
			return "", fmt.Errorf("submission failed with error type `%s` and description `%s`", failure.Type, failure.Description)
		}
		return "", fmt.Errorf("an unknown submission failed: %v", resp.NotCreated)
	}
	for k, v := range resp.Created {
		if k == requestID.String() {
			return v.ID, nil
		}
	}
	return "", fmt.Errorf("request id %s not found", requestID.String())
}

func getIdentityID(c *client.Client, email string) (string, error) {
	// TODO: retrieve identity id
	return "", nil
}
