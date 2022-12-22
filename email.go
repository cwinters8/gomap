package gomap

import (
	"fmt"
)

type Email struct {
	From    string
	To      []string
	Subject string
	Message string
	mailbox *Mailbox
}

func (m *Mailbox) NewEmail(from string, to []string) *Email {
	return &Email{
		mailbox: m,
		From:    from,
		To:      to,
	}
}

func (e *Email) Send() error {
	return fmt.Errorf("not implemented")
}
