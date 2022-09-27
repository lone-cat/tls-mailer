package simpleemail

import (
	"errors"
	"github.com/lone-cat/tls-mailer/simpleemail/address"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"github.com/lone-cat/tls-mailer/simpleemail/part"
	"net/mail"
	"strings"
)

type Email struct {
	headers headers.Headers

	from address.AddressList
	to   address.AddressList
	cc   address.AddressList
	bcc  address.AddressList

	subject []byte

	mainPart *relatedSubPart

	attachments part.PartsList
}

func NewEmptyEmail() *Email {
	return &Email{
		headers: headers.NewHeaders(),

		from: address.NewAddressList(),
		to:   address.NewAddressList(),
		cc:   address.NewAddressList(),
		bcc:  address.NewAddressList(),

		subject: make([]byte, 0),

		mainPart: newRelatedSubPart(),

		attachments: part.NewPartsList(),
	}
}

func (e *Email) GetFrom() []*mail.Address {
	return e.from.ExportMailAddressSlice()
}

func (e *Email) GetTo() []*mail.Address {
	return e.to.ExportMailAddressSlice()
}

func (e *Email) GetCc() []*mail.Address {
	return e.cc.ExportMailAddressSlice()
}

func (e *Email) GetBcc() []*mail.Address {
	return e.bcc.ExportMailAddressSlice()
}

func (e *Email) GetSubject() string {
	return string(e.subject)
}

func (e *Email) GetText() string {
	return string(e.mainPart.alternativeSubPart.textPart.GetBody())
}

func (e *Email) GetHtml() string {
	return string(e.mainPart.alternativeSubPart.htmlPart.GetBody())
}

func (e *Email) WithFrom(from []*mail.Address) *Email {
	newEmail := e.clone()
	newEmail.from = newEmail.from.WithMailAddressSlice(from...)
	return newEmail
}

func (e *Email) WithTo(to []*mail.Address) *Email {
	newEmail := e.clone()
	newEmail.to = newEmail.to.WithMailAddressSlice(to...)
	return newEmail
}

func (e *Email) WithCc(cc []*mail.Address) *Email {
	newEmail := e.clone()
	newEmail.cc = newEmail.cc.WithMailAddressSlice(cc...)
	return newEmail
}

func (e *Email) WithBcc(bcc []*mail.Address) *Email {
	newEmail := e.clone()
	newEmail.bcc = newEmail.bcc.WithMailAddressSlice(bcc...)
	return newEmail
}

func (e *Email) WithSubject(subject string) *Email {
	newEmail := e.clone()
	newEmail.subject = []byte(subject)
	return newEmail
}

func (e *Email) WithText(text string) *Email {
	newEmail := e.clone()
	newEmail.mainPart = newEmail.mainPart.WithText([]byte(text))
	return newEmail
}

func (e *Email) WithHtml(html string) *Email {
	newEmail := e.clone()
	newEmail.mainPart = newEmail.mainPart.WithHtml([]byte(html))
	return newEmail
}

func (e *Email) WithEmbeddedFile(cid string, filename string) (*Email, error) {
	embedded, err := part.NewEmbeddedPartFromFile(cid, filename)
	if err != nil {
		return e, err
	}
	newEmail := e.withEmbedded(embedded)
	return newEmail, nil
}

func (e *Email) WithEmbeddedBytes(cid string, bts []byte) *Email {
	embedded := part.NewEmbeddedPartFromBytes(cid, bts)
	return e.withEmbedded(embedded)
}

func (e *Email) withEmbedded(embedded part.Part) *Email {
	newEmail := e.clone()
	newEmail.mainPart = newEmail.mainPart.WithEmbeddedSubPart(embedded)
	return newEmail
}

func (e *Email) WithoutEmbedded() *Email {
	newEmail := e.clone()
	newEmail.mainPart = newEmail.mainPart.WithoutEmbeddedSubParts()
	return newEmail
}

func (e *Email) WithAttachedFile(filename string) (*Email, error) {
	attachment, err := part.NewAttachedPartFromFile(filename)
	if err != nil {
		return e, err
	}
	newEmail := e.withAttached(attachment)
	return newEmail, nil
}

func (e *Email) WithAttachedBytes(bts []byte) *Email {
	attachment := part.NewAttachedPartFromBytes(bts)
	return e.withAttached(attachment)
}

func (e *Email) withAttached(attachment part.Part) *Email {
	newEmail := e.clone()
	newEmail.attachments = newEmail.attachments.WithAppended(attachment)
	return newEmail
}

func (e *Email) WithoutAttachments() *Email {
	newEmail := e.clone()
	newEmail.attachments = part.NewPartsList()
	return newEmail
}

func (e *Email) Compile() ([]byte, error) {
	exportedPart := e.toPart()

	eFrom := e.from.ExportMailAddressSlice()
	if len(eFrom) > 0 {
		from := make([]string, len(eFrom))
		for index, addr := range eFrom {
			from[index] = addr.String()
		}
		froms := strings.Join(from, `, `)
		err := validateAddresses(froms)
		if err != nil {
			return nil, err
		}
		exportedPart = exportedPart.WithHeaders(exportedPart.GetHeaders().WithHeader(`from`, froms))
	}

	eTo := e.to.ExportMailAddressSlice()
	rcpts := make([]string, 0)
	if len(eTo) > 0 {
		to := make([]string, len(eTo))
		for index, addr := range eTo {
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

	eCc := e.cc.ExportMailAddressSlice()
	if len(eCc) > 0 {
		cc := make([]string, len(eCc))
		for index, addr := range eCc {
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

	eBcc := e.bcc.ExportMailAddressSlice()
	if len(eBcc) > 0 {
		bcc := make([]string, len(eBcc))
		for index, addr := range eBcc {
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
	if len(e.subject) > 0 {
		exportedPart = exportedPart.WithHeaders(exportedPart.GetHeaders().WithHeader(`subject`, string(e.subject)))
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
	eFrom := e.from.ExportMailAddressSlice()
	if len(eFrom) < 1 {
		return nil
	}
	return &mail.Address{Name: eFrom[0].Name, Address: eFrom[0].Address}
}

func (e *Email) GetRecipients() []*mail.Address {
	eTo := e.to.ExportMailAddressSlice()
	recipients := make([]*mail.Address, len(eTo))
	for ind, recipient := range eTo {
		recipients[ind] = &mail.Address{Name: recipient.Name, Address: recipient.Address}
	}
	return recipients
}

func (e *Email) clone() *Email {
	cloneEmail := *e
	return &cloneEmail
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
