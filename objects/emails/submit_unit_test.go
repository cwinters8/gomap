package emails_test

import (
	"testing"

	"github.com/cwinters8/gomap/objects/emails"
	"github.com/google/uuid"
)

func TestSubmitCallUsesProvidedIdentity(t *testing.T) {
	requestID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	identityID := "identity-owned-by-account"
	emailID := "email-with-custom-from-header"
	call, err := emails.SubmitCall(requestID, identityID, "account-id", emailID, "drafts-id", "sent-id")
	if err != nil {
		t.Fatalf("failed to construct submit call: %v", err)
	}
	if call.Method != "EmailSubmission/set" {
		t.Fatalf("wanted EmailSubmission/set; got %s", call.Method)
	}
	create, ok := call.Arguments["create"].(map[string]map[string]string)
	if !ok {
		t.Fatalf("create argument had unexpected type %T", call.Arguments["create"])
	}
	created := create[requestID.String()]
	if created["identityId"] != identityID {
		t.Fatalf("wanted identityId %s; got %s", identityID, created["identityId"])
	}
	if created["emailId"] != emailID {
		t.Fatalf("wanted emailId %s; got %s", emailID, created["emailId"])
	}
}

func TestSubmitRequiresFromAddressForDefaultIdentity(t *testing.T) {
	email := &emails.Email{}
	_, err := email.Submit(nil, "drafts-id", "sent-id")
	if err == nil {
		t.Fatalf("wanted error for missing from address")
	}
}
