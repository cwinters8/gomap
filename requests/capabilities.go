package requests

const (
	UsingCore       Capability = "urn:ietf:params:jmap:core"
	UsingMail       Capability = "urn:ietf:params:jmap:mail"
	UsingSubmission Capability = "urn:ietf:params:jmap:submission"
)

type Capability string
