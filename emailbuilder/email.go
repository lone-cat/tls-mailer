package emailbuilder

import (
	"net/mail"
	"strings"
)

type Email struct {
	headers Headers

	From []mail.Address
	To   []mail.Address
	Cc   []mail.Address
	Bcc  []mail.Address

	Subject string

	Html string
	Text string

	embedded    subParts
	attachments subParts
}

func newEmail() Email {
	return Email{
		headers: NewHeaders(),

		From: make([]mail.Address, 0),
		To:   make([]mail.Address, 0),
		Cc:   make([]mail.Address, 0),
		Bcc:  make([]mail.Address, 0),

		embedded:    NewSubParts(),
		attachments: NewSubParts(),
	}
}

func (e Email) GetHeaders() Headers {
	return e.headers.Clone()
}

func (e Email) Render1() string {
	headersToRender := e.GetHeaders()
	headersToRender = headersToRender.WithHeader(SubjectHeader, e.Subject)
	if len(e.From) > 0 {
		headersToRender = headersToRender.WithHeader(FromHeader, e.From[0].String())
	}
	if len(e.To) > 0 {
		toStrings := make([]string, len(e.To))
		for index, to := range e.To {
			toStrings[index] = to.String()
		}
		headersToRender = headersToRender.WithHeader(ToHeader, strings.Join(toStrings, `, `))
	}

	return headersToRender.Render()
}

func (e Email) Render() (string, error) {
	part := newPart()
	part.headers = e.headers.Clone()
	part.subParts = append(e.embedded, e.attachments...)
	part.subParts = append(part.subParts, part.subParts...)

	return part.Render()
}

///
///
///
func (e Email) GetAttachments() int {
	return len(e.attachments)
}
