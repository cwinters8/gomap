package gomap

import (
	"github.com/cwinters8/gomap/mail"
)

type Object interface {
	*mail.Mailbox
	MethodPrefix() string
	Query() error
}
