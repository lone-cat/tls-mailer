package simpleemail

import (
	"errors"
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail/address"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"github.com/lone-cat/tls-mailer/simpleemail/part"
	"net/http"
	"net/mail"
	"strings"
)

type Email interface {
	fmt.Stringer
	GetFrom() []*mail.Address
	GetTo() []*mail.Address
	GetCc() []*mail.Address
	GetBcc() []*mail.Address
	GetSubject() string
	GetText() string
	GetHtml() string
	Compile() ([]byte, error)
	WithFrom(...*mail.Address) (Email, error)
	WithTo(...*mail.Address) (Email, error)
	WithCc(...*mail.Address) (Email, error)
	WithBcc(...*mail.Address) (Email, error)
	WithSubject(string) Email
	WithText(string) (Email, error)
	WithHtml(string) (Email, error)
	WithEmbeddedFile(cid string, filename string) (Email, error)
	WithEmbeddedBytes(cid string, bts []byte) Email
	WithoutEmbedded() Email
	WithAttachedFile(filename string) (Email, error)
	WithAttachedBytes(bts []byte) Email
	WithoutAttachments() Email
	Dump() map[string]any
}

type email struct {
	headers headers.Headers

	from address.AddressList
	to   address.AddressList
	cc   address.AddressList
	bcc  address.AddressList

	subject []byte

	mainPart *relatedSubPart

	attachments part.PartsList
}

func NewEmptyEmail() Email {
	return &email{
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

func (e *email) GetFrom() []*mail.Address {
	return e.from.ExportMailAddressSlice()
}

func (e *email) GetTo() []*mail.Address {
	return e.to.ExportMailAddressSlice()
}

func (e *email) GetCc() []*mail.Address {
	return e.cc.ExportMailAddressSlice()
}

func (e *email) GetBcc() []*mail.Address {
	return e.bcc.ExportMailAddressSlice()
}

func (e *email) GetSubject() string {
	return string(e.subject)
}

func (e *email) GetText() string {
	return string(e.mainPart.alternativeSubPart.textPart.GetBody())
}

func (e *email) GetHtml() string {
	return string(e.mainPart.alternativeSubPart.htmlPart.GetBody())
}

func (e *email) WithFrom(from ...*mail.Address) (mail Email, err error) {
	if err = validateMailAddressSlice(from); err != nil {
		return
	}
	newEmail := e.clone()
	newEmail.from = newEmail.from.WithMailAddressSlice(from...)
	mail = newEmail
	return
}

func (e *email) WithTo(to ...*mail.Address) (mail Email, err error) {
	if err = validateMailAddressSlice(to); err != nil {
		return
	}
	newEmail := e.clone()
	newEmail.to = newEmail.to.WithMailAddressSlice(to...)
	mail = newEmail
	return
}

func (e *email) WithCc(cc ...*mail.Address) (mail Email, err error) {
	if err = validateMailAddressSlice(cc); err != nil {
		return
	}
	newEmail := e.clone()
	newEmail.cc = newEmail.cc.WithMailAddressSlice(cc...)
	mail = newEmail
	return
}

func (e *email) WithBcc(bcc ...*mail.Address) (mail Email, err error) {
	if err = validateMailAddressSlice(bcc); err != nil {
		return
	}
	newEmail := e.clone()
	newEmail.bcc = newEmail.bcc.WithMailAddressSlice(bcc...)
	mail = newEmail
	return
}

func (e *email) WithSubject(subject string) Email {
	newEmail := e.clone()
	newEmail.subject = []byte(subject)
	return newEmail
}

func (e *email) WithText(text string) (mail Email, err error) {
	textBts := []byte(text)
	ct := http.DetectContentType(textBts)
	if !strings.HasPrefix(ct, headers.TextPlain) {
		err = fmt.Errorf(`content type "%s" passed as text part`, ct)
	}
	newEmail := e.clone()
	newEmail.mainPart = newEmail.mainPart.WithText(textBts)
	mail = newEmail
	return
}

func (e *email) WithHtml(html string) (mail Email, err error) {
	htmlBts := []byte(html)
	ct := http.DetectContentType(htmlBts)
	if !strings.HasPrefix(ct, headers.TextHtml) {
		err = fmt.Errorf(`content type "%s" passed as html part`, ct)
	}
	newEmail := e.clone()
	newEmail.mainPart = newEmail.mainPart.WithHtml(htmlBts)
	mail = newEmail
	return
}

func (e *email) WithEmbeddedFile(cid string, filename string) (Email, error) {
	embedded, err := part.NewEmbeddedPartFromFile(cid, filename)
	if err != nil {
		return e, err
	}
	newEmail := e.withEmbedded(embedded)
	return newEmail, nil
}

func (e *email) WithEmbeddedBytes(cid string, bts []byte) Email {
	embedded := part.NewEmbeddedPartFromBytes(cid, bts)
	return e.withEmbedded(embedded)
}

func (e *email) withEmbedded(embedded part.Part) Email {
	newEmail := e.clone()
	newEmail.mainPart = newEmail.mainPart.WithEmbeddedSubPart(embedded)
	return newEmail
}

func (e *email) WithoutEmbedded() Email {
	newEmail := e.clone()
	newEmail.mainPart = newEmail.mainPart.WithoutEmbeddedSubParts()
	return newEmail
}

func (e *email) WithAttachedFile(filename string) (Email, error) {
	attachment, err := part.NewAttachedPartFromFile(filename)
	if err != nil {
		return e, err
	}
	newEmail := e.withAttached(attachment)
	return newEmail, nil
}

func (e *email) WithAttachedBytes(bts []byte) Email {
	attachment := part.NewAttachedPartFromBytes(bts)
	return e.withAttached(attachment)
}

func (e *email) withAttached(attachment part.Part) Email {
	newEmail := e.clone()
	newEmail.attachments = newEmail.attachments.WithAppended(attachment)
	return newEmail
}

func (e *email) WithoutAttachments() Email {
	newEmail := e.clone()
	newEmail.attachments = part.NewPartsList()
	return newEmail
}

func (e *email) Compile() ([]byte, error) {
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

func (e *email) String() string {
	compiled, err := e.Compile()
	if err != nil {
		return err.Error()
	}
	return string(compiled)
}

func (e *email) GetSender() *mail.Address {
	eFrom := e.GetFrom()
	if len(eFrom) < 1 {
		return nil
	}
	return eFrom[0]
}

func (e *email) GetRecipients() []*mail.Address {
	return e.GetFrom()
}

func (e *email) clone() *email {
	cloneEmail := *e
	return &cloneEmail
}

func (e *email) toPart() part.Part {
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

func (e *email) Dump() map[string]any {
	if e == nil {
		return nil
	}

	dump := make(map[string]any)
	dump[`to`] = e.to.Dump()
	dump[`from`] = e.from.Dump()
	dump[`cc`] = e.cc.Dump()
	dump[`bcc`] = e.bcc.Dump()
	dump[`relatedPart`] = e.mainPart.Dump()
	dump[`attachedParts`] = e.attachments.Dump()

	dump[`compiled`] = func() (compiled string) {
		defer func() {
			if r := recover(); r != nil {
				if compiled != `` {
					compiled += "\r\n"
				}
				compiled = compiled + fmt.Sprintf("%#v", r)
			}
		}()
		compiled = e.String()
		return
	}()

	return dump
}

func validateAddresses(addrs string) error {
	_, err := mail.ParseAddressList(addrs)
	if err == nil {
		return err
	}
	return errors.New(`address contains invalid value`)
}

func validateMailAddressSlice(addrss []*mail.Address) (err error) {
	for i := range addrss {
		if addrss[i] == nil {
			err = errors.New(`nil *mail.Address passed`)
			return
		}
		err = validateEmail(addrss[i].Address)
		if err != nil {
			return
		}
	}

	return
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New(`invalid email address "` + email + `"`)
	}

	return nil
}
