[![Go Reference](https://pkg.go.dev/badge/github.com/cwinters8/gomap.svg)](https://pkg.go.dev/github.com/cwinters8/gomap) [![Go Report Card](https://goreportcard.com/badge/github.com/cwinters8/gomap)](https://goreportcard.com/report/github.com/cwinters8/gomap) [![License](https://img.shields.io/github/license/cwinters8/gomap?color=blue)](https://github.com/cwinters8/gomap/blob/main/LICENSE) <!-- markdownlint-disable-line MD041 -->

# gomap

Go module for interfacing with [JMAP](https://jmap.io) mail servers, such as [Fastmail](https://www.fastmail.com/).

## Usage

To add the module to your project:

```sh
go get -u github.com/cwinters8/gomap
```

First, create a new mail client. You will need a session URL and a bearer token from your JMAP mail server. For Fastmail, the session URL is `https://api.fastmail.com/jmap/session`, and you can create an API token in [Settings > Privacy & Security > Integrations > API tokens](https://www.fastmail.com/settings/security/tokens). You will most likely need to give access to `Email` (for querying and reading email contents) and `Email submission` (for sending emails) when you create the token.

```go
mail, err := gomap.NewClient(
  "https://api.fastmail.com/jmap/session",
  os.Getenv("BEARER_TOKEN"),
  gomap.DefaultDrafts,
  gomap.DefaultSent,
)
if err != nil {
  log.Fatal(err)
}
```

Then you can use the client for your chosen operations. Check out the [examples](https://pkg.go.dev/github.com/cwinters8/gomap#pkg-examples) for full details on how to send and find emails.

### Sending with a custom visible From address

JMAP requires every `EmailSubmission` to be associated with an account `Identity`,
but the message's visible `From` header is part of the `Email` object. Use
`SendEmailWithIdentity` when those values need to differ, such as contact-form
notifications where the visible sender should be the email address submitted in
the form while the authenticated mailbox identity is used for submission:

```go
submitter := gomap.NewAddress("Form Submitter", "submitter@example.com")
recipient := gomap.NewAddress("Contact Inbox", "contact@example.org")

err := mail.SendEmailWithIdentity(
  gomap.NewAddresses(submitter),
  gomap.NewAddresses(recipient),
  "Contact form submission",
  "Message body",
  "administrator@example.org", // authenticated JMAP identity email
  false,
)
```

Servers may still reject this if the selected identity is not permitted to use
the requested `From` header.
