package gomap_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cwinters8/gomap"

	"github.com/joho/godotenv"
)

func TestNewClient(t *testing.T) {
	env(t)
	client, err := gomap.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		failf(t, "failed to instantiate new client: %s", err.Error())
	}
	wantURL := "https://api.fastmail.com/jmap/api/"
	switch {
	case client == nil:
		failf(t, "returned client should not be nil")
	case client.Session == nil:
		t.Error("session should not be nil")
	case client.HTTPClient == nil:
		t.Error("http client should not be nil")
	case client.Session.APIURL != wantURL:
		t.Errorf("wanted URL %s; got URL %s", wantURL, client.Session.APIURL)
	}
}

func env(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Logf("failed to load .env: %s", err.Error())
		t.FailNow()
	}
}

// logs a formatted string, then fails the test immediately
func failf(t *testing.T, format string, args ...any) {
	t.Logf(fmt.Sprintf("%s\n", format), args...)
	t.FailNow()
}
