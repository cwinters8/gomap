package gomap_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cwinters8/gomap"
	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"
)

func ExampleClient_SendEmail() {
	mail, err := gomap.NewClient(
		"https://api.fastmail.com/jmap/session",
		os.Getenv("BEARER_TOKEN"),
		gomap.DefaultDrafts,
		gomap.DefaultSent,
	)
	if err != nil {
		log.Fatal(err)
	}

	// send an email
	from := gomap.NewAddress("Clark Winters", "dev@clarkwinters.com")
	to := gomap.NewAddress("Tester Gopher", "tester@clarkwinters.com")

	if err := mail.SendEmail(
		gomap.NewAddresses(from),
		gomap.NewAddresses(to),
		"Hello from gomap!",
		"Hello Tester Gopher,\n\nNice to meet you.",
		false,
	); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_GetEmails() {
	mail, err := gomap.NewClient(
		"https://api.fastmail.com/jmap/session",
		os.Getenv("BEARER_TOKEN"),
		gomap.DefaultDrafts,
		gomap.DefaultSent,
	)
	if err != nil {
		log.Fatal(err)
	}

	// send an email with a unique identifier
	from := gomap.NewAddress("Clark Winters", "dev@clarkwinters.com")
	to := gomap.NewAddress("Tester Gopher", "tester@clarkwinters.com")

	id := uuid.New()
	strID := id.String()

	if err := mail.SendEmail(
		gomap.NewAddresses(from),
		gomap.NewAddresses(to),
		"Hello from gomap!",
		fmt.Sprintf("ID: %s", strID),
		false,
	); err != nil {
		log.Fatal(err)
	}

	// retrieve the email
	emails, err := mail.GetEmails(&gomap.Filter{Text: strID}, 1, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	if len(emails) > 0 {
		email := emails[0]
		fmt.Println("found email:")
		fmt.Printf("\tfrom: %+v\n", *email.From[0])
		fmt.Printf("\tto: %+v\n", *email.To[0])
		fmt.Printf("\tsubject: %s\n", email.Subject)
		fmt.Printf("\tbody: %s\n", email.Body.Value)
	}

	// Output:
	// found email:
	//	from: {Name:Clark Winters Email:dev@clarkwinters.com}
	//	to: {Name:Tester Gopher Email:tester@clarkwinters.com}
	//	subject: Hello from gomap!
	//	body: ID: 68784752-95e1-4fc6-b923-0e84aafe1150
}

func TestGetEmails(t *testing.T) {
	if err := utils.Env(".env"); err != nil {
		t.Fatalf("failed to load env variables: %v", err)
	}
	ExampleClient_GetEmails()
}
