package requests_test

import (
	"testing"

	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"
)

func TestNewGet(t *testing.T) {
	t.Run("empty IDs and properties", func(t *testing.T) {
		want := requests.Get{
			Prefix: "Mailbox",
			Body:   &requests.GetBody{AccountID: "xyz"},
		}
		got, err := requests.NewGet(want.Body.AccountID, want.Prefix, nil, nil)
		if err != nil {
			t.Fatalf("failed to instantiate new Get: %s", err.Error())
		}
		cases := utils.Cases{
			utils.NewCase(
				got.Prefix != want.Prefix,
				"wanted prefix %s; got %s",
				want.Prefix, got.Prefix,
			),
			utils.NewCase(
				got.Body.AccountID != want.Body.AccountID,
				"wanted account id %s; got %s",
				want.Body.AccountID, got.Body.AccountID,
			),
			utils.NewCase(
				len(got.Body.IDs) != 0,
				"wanted IDs slice to be empty; got %d",
				len(got.Body.IDs),
			),
			utils.NewCase(
				len(got.Body.Properties) != 0,
				"wanted Properties slice to be empty; got %d",
				len(got.Body.Properties),
			),
		}
		cases.Iterator(func(c *utils.Case) {
			t.Error(c.Message)
		})
	})
	t.Run("non-empty IDs and properties", func(t *testing.T) {
		wantID := "xyz-id"
		wantProp := "name"
		want := requests.Get{
			Prefix: "Mailbox",
			Body: &requests.GetBody{
				AccountID:  "xyz",
				IDs:        []string{wantID},
				Properties: []string{wantProp},
			},
		}
		got, err := requests.NewGet(want.Body.AccountID, want.Prefix, want.Body.IDs, want.Body.Properties)
		if err != nil {
			t.Fatalf("failed to instantiate new Get: %s", err.Error())
		}
		cases := utils.Cases{
			utils.NewCase(
				got.Body.IDs[0] != wantID,
				"wanted body id %s; got %s",
				wantID, got.Body.IDs[0],
			),
			utils.NewCase(
				got.Body.Properties[0] != wantProp,
				"wanted body property %s; got %s",
				wantProp, got.Body.Properties[0],
			),
		}
		cases.Iterator(func(c *utils.Case) {
			t.Error(c.Message)
		})
	})
}
