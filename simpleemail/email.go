package simpleemail

import (
	"errors"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"github.com/lone-cat/tls-mailer/simpleemail/part"
	"net/mail"
	"strings"
)

type Email struct {
	headers headers.Headers

	from addresses
	to   addresses
	cc   addresses
	bcc  addresses

	subject string

	mainPart *relatedSubPart

	attachments part.subParts
}

func NewEmptyEmail() *Email {
	return &Email{
		headers: headers.NewHeaders(),

		from: newAddresses(),
		to:   newAddresses(),
		cc:   newAddresses(),
		bcc:  newAddresses(),

		mainPart: newRelatedSubPart(),

		attachments: part.newSubParts(),
	}
}

func (e *Email) GetFrom() []*mail.Address {
	return e.from.clone()
}

func (e *Email) GetTo() []*mail.Address {
	return e.to.clone()
}

func (e *Email) GetCc() []*mail.Address {
	return e.cc.clone()
}

func (e *Email) GetBcc() []*mail.Address {
	return e.bcc.clone()
}

func (e *Email) GetSubject() string {
	return e.subject
}

func (e *Email) GetText() string {
	return e.mainPart.alternativeSubPart.textPart.body
}

func (e *Email) GetHtml() string {
	return e.mainPart.alternativeSubPart.htmlPart.body
}

func (e *Email) WithFrom(from []*mail.Address) *Email {
	newEmail := e.clone()
	newEmail.from = from
	return newEmail
}

func (e *Email) WithTo(to []*mail.Address) *Email {
	newEmail := e.clone()
	newEmail.to = to
	return newEmail
}

func (e *Email) WithCc(cc []*mail.Address) *Email {
	newEmail := e.clone()
	newEmail.cc = cc
	return newEmail
}

func (e *Email) WithBcc(bcc []*mail.Address) *Email {
	newEmail := e.clone()
	newEmail.bcc = bcc
	return newEmail
}

func (e *Email) WithSubject(subject string) *Email {
	newEmail := e.clone()
	newEmail.subject = subject
	return newEmail
}

func (e *Email) WithText(text string) *Email {
	newEmail := e.clone()
	newEmail.mainPart.alternativeSubPart = newEmail.mainPart.alternativeSubPart.withText(text)
	return newEmail
}

func (e *Email) WithHtml(html string) *Email {
	newEmail := e.clone()
	newEmail.mainPart.alternativeSubPart = newEmail.mainPart.alternativeSubPart.withHtml(html)
	return newEmail
}

func (e *Email) WithEmbeddedFile(cid string, filename string) (*Email, error) {
	embedded, err := part.newEmbeddedPartFromFile(cid, filename)
	if err != nil {
		return e, err
	}
	newEmail := e.clone()
	newEmail.mainPart.embeddedSubParts = append(newEmail.mainPart.embeddedSubParts, embedded)
	return newEmail, nil
}

func (e *Email) WithoutEmbeddedFiles() *Email {
	newEmail := e.clone()
	newEmail.mainPart.embeddedSubParts = part.newSubParts()
	return newEmail
}

func (e *Email) WithAttachedFile(filename string) (*Email, error) {
	attachment, err := part.newAttachedPartFromFile(filename)
	if err != nil {
		return e, err
	}
	newEmail := e.clone()
	newEmail.attachments = append(newEmail.attachments, attachment)
	return newEmail, nil
}

func (e *Email) WithoutAttachedFiles() *Email {
	newEmail := e.clone()
	newEmail.attachments = part.newSubParts()
	return newEmail
}

func (e *Email) Compile() ([]byte, error) {
	exportedPart := e.toPart()

	if len(e.from) > 0 {
		from := make([]string, len(e.from))
		for index, addr := range e.from {
			from[index] = addr.String()
		}
		froms := strings.Join(from, `, `)
		err := validateAddresses(froms)
		if err != nil {
			return nil, err
		}
		exportedPart.headers = exportedPart.headers.WithHeader(`from`, froms)
	}

	rcpts := make([]string, 0)
	if len(e.to) > 0 {
		to := make([]string, len(e.to))
		for index, addr := range e.to {
			to[index] = addr.String()
			rcpts = append(rcpts, addr.Address)
		}
		tos := strings.Join(to, `, `)
		err := validateAddresses(tos)
		if err != nil {
			return nil, err
		}
		exportedPart.headers = exportedPart.headers.WithHeader(`to`, tos)
	}
	if len(e.cc) > 0 {
		cc := make([]string, len(e.cc))
		for index, addr := range e.cc {
			cc[index] = addr.String()
			rcpts = append(rcpts, addr.Address)
		}
		ccs := strings.Join(cc, `, `)
		err := validateAddresses(ccs)
		if err != nil {
			return nil, err
		}
		exportedPart.headers = exportedPart.headers.WithHeader(`cc`, ccs)
	}
	if len(e.bcc) > 0 {
		bcc := make([]string, len(e.bcc))
		for index, addr := range e.bcc {
			bcc[index] = addr.String()
			rcpts = append(rcpts, addr.Address)
		}
		bccs := strings.Join(bcc, `, `)
		err := validateAddresses(bccs)
		if err != nil {
			return nil, err
		}
		exportedPart.headers = exportedPart.headers.WithHeader(`bcc`, bccs)
	}
	if e.subject != `` {
		exportedPart.headers = exportedPart.headers.WithHeader(`subject`, e.subject)
	}

	return exportedPart.compile()
}

func (e *Email) String() string {
	compiled, err := e.Compile()
	if err != nil {
		return err.Error()
	}
	return string(compiled)
}

func (e *Email) GetSender() *mail.Address {
	if len(e.from) < 1 {
		return nil
	}
	return &mail.Address{Name: e.from[0].Name, Address: e.from[0].Address}
}

func (e *Email) GetRecipients() []*mail.Address {
	recipients := make([]*mail.Address, len(e.to))
	for ind, recipient := range e.to {
		recipients[ind] = &mail.Address{Name: recipient.Name, Address: recipient.Address}
	}
	return recipients
}

func (e *Email) clone() *Email {

	newEmail := NewEmptyEmail()

	newEmail.headers = e.headers.clone()

	newEmail.from = e.from.clone()
	newEmail.to = e.to.clone()
	newEmail.cc = e.cc.clone()
	newEmail.bcc = e.bcc.clone()

	newEmail.subject = e.subject

	newEmail.mainPart = e.mainPart.clone()
	newEmail.attachments = e.attachments.clone()

	return newEmail
}

func (e *Email) toPart() *part.part {
	mainPart := e.mainPart.toPart()

	if len(e.attachments) < 1 {
		return mainPart
	}

	exportedPart := &part.part{
		headers:  e.headers.clone(),
		subParts: append([]*part.part{mainPart}, e.attachments...),
	}

	if !exportedPart.headers.IsMultipart() {
		exportedPart.headers = exportedPart.headers.WithHeader(`Content-Type`, part.MultipartMixed)
	}

	return exportedPart
}

func validateAddresses(addrs string) error {
	_, err := mail.ParseAddressList(addrs)
	if err == nil {
		return err
	}
	return errors.New(`address contains invalid value`)
}
