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

	attachments part.PartsList
}

func NewEmptyEmail() *Email {
	return &Email{
		headers: headers.NewHeaders(),

		from: newAddresses(),
		to:   newAddresses(),
		cc:   newAddresses(),
		bcc:  newAddresses(),

		mainPart: newRelatedSubPart(),

		attachments: part.NewPartsList(),
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
	return e.mainPart.alternativeSubPart.textPart.GetBody()
}

func (e *Email) GetHtml() string {
	return e.mainPart.alternativeSubPart.htmlPart.GetBody()
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
	embedded, err := part.NewEmbeddedPartFromFile(cid, filename)
	if err != nil {
		return e, err
	}
	newEmail := e.clone()
	newEmail.mainPart.embeddedSubParts = newEmail.mainPart.embeddedSubParts.WithAppended(embedded)
	return newEmail, nil
}

func (e *Email) WithoutEmbeddedFiles() *Email {
	newEmail := e.clone()
	newEmail.mainPart.embeddedSubParts = part.NewPartsList()
	return newEmail
}

func (e *Email) WithAttachedFile(filename string) (*Email, error) {
	attachment, err := part.NewAttachedPartFromFile(filename)
	if err != nil {
		return e, err
	}
	newEmail := e.clone()
	newEmail.attachments = newEmail.attachments.WithAppended(attachment)
	return newEmail, nil
}

func (e *Email) WithoutAttachedFiles() *Email {
	newEmail := e.clone()
	newEmail.attachments = part.NewPartsList()
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
		exportedPart = exportedPart.WithHeaders(exportedPart.GetHeaders().WithHeader(`from`, froms))
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
		exportedPart = exportedPart.WithHeaders(exportedPart.GetHeaders().WithHeader(`to`, tos))
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
		exportedPart = exportedPart.WithHeaders(exportedPart.GetHeaders().WithHeader(`cc`, ccs))
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
		exportedPart = exportedPart.WithHeaders(exportedPart.GetHeaders().WithHeader(`bcc`, bccs))
	}
	if e.subject != `` {
		exportedPart = exportedPart.WithHeaders(exportedPart.GetHeaders().WithHeader(`subject`, e.subject))
	}

	return exportedPart.Compile()
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

	newEmail.headers = e.headers.Clone()

	newEmail.from = e.from.clone()
	newEmail.to = e.to.clone()
	newEmail.cc = e.cc.clone()
	newEmail.bcc = e.bcc.clone()

	newEmail.subject = e.subject

	newEmail.mainPart = e.mainPart.clone()
	newEmail.attachments = part.NewPartsList(e.attachments.ExtractPartsSlice()...)

	return newEmail
}

func (e *Email) toPart() part.Part {
	mainPart := e.mainPart.toPart()

	if len(e.attachments.ExtractPartsSlice()) < 1 {
		return mainPart
	}

	exportedPart := part.NewPart().WithHeaders(e.headers).WithSubParts(append([]part.Part{mainPart}, e.attachments.ExtractPartsSlice()...)...)

	if !exportedPart.GetHeaders().IsMultipart() {
		exportedPart = exportedPart.WithHeaders(exportedPart.GetHeaders().WithHeader(`Content-Type`, headers.MultipartMixed))
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
