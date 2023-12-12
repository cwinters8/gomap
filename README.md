# gomap

[![Go Reference](https://pkg.go.dev/badge/github.com/cwinters8/gomap.svg)](https://pkg.go.dev/github.com/cwinters8/gomap)

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
