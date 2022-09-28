package part

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/common"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type part struct {
	headers  headers.Headers
	body     []byte
	subParts PartsList
}

func NewPart() Part {
	return &part{
		headers:  headers.NewHeaders(),
		subParts: NewPartsList(),
	}
}

func NewPartFromString(data string) (createdPart Part) {
	return NewPart().WithBodyFromString(data)
}

func NewPartFromBytes(data []byte) (createdPart Part) {
	return NewPart().WithBody(data)
}

func NewEmbeddedPartFromString(cid string, data string) (embeddedPart Part) {
	embeddedPart = NewPartFromString(data)
	hdrs := embeddedPart.GetHeaders().
		WithHeader(
			headers.ContentDisposition,
			`inline`,
		)
	if cid != `` {
		hdrs = hdrs.WithHeader(
			headers.ContentId,
			fmt.Sprintf(`<%s>`, cid),
		)
	}

	embeddedPart = embeddedPart.WithHeaders(hdrs)

	return
}

func NewEmbeddedPartFromBytes(cid string, data []byte) (embeddedPart Part) {
	embeddedPart = NewPartFromBytes(data)
	hdrs := embeddedPart.GetHeaders().
		WithHeader(
			headers.ContentDisposition,
			`inline`,
		)
	if cid != `` {
		hdrs = hdrs.WithHeader(
			headers.ContentId,
			fmt.Sprintf(`<%s>`, cid),
		)
	}

	embeddedPart = embeddedPart.WithHeaders(hdrs)

	return
}

func NewEmbeddedPartFromFile(cid, filename string) (embedded Part, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	embedded = NewEmbeddedPartFromBytes(cid, data)
	hdrs := embedded.GetHeaders().
		WithHeader(
			headers.ContentDisposition,
			fmt.Sprintf(`inline; filename="%s"`, filepath.Base(filename)),
		)

	embedded = embedded.WithHeaders(hdrs)

	return
}

func NewAttachedPartFromString(data string) (attachment Part) {
	attachment = NewPartFromString(data)
	hdrs := attachment.GetHeaders().
		WithHeader(
			headers.ContentDisposition,
			`attachment`,
		)

	attachment = attachment.WithHeaders(hdrs)

	return
}

func NewAttachedPartFromBytes(data []byte) (attachment Part) {
	attachment = NewPartFromBytes(data)
	hdrs := attachment.GetHeaders().
		WithHeader(
			headers.ContentDisposition,
			`attachment`,
		)

	attachment = attachment.WithHeaders(hdrs)

	return
}

func NewAttachedPartFromFile(filename string) (attachment Part, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	attachment = NewAttachedPartFromBytes(data)
	hdrs := attachment.GetHeaders().
		WithHeader(
			headers.ContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filename)),
		)

	attachment = attachment.WithHeaders(hdrs)

	return
}

func (p *part) Clone() Part {
	return &part{
		headers:  p.headers,
		body:     p.body,
		subParts: p.subParts,
	}
}

func (p *part) GetHeaders() headers.Headers {
	return p.headers
}

func (p *part) WithHeaders(headers headers.Headers) Part {
	return &part{
		headers:  headers,
		body:     p.body,
		subParts: p.subParts,
	}
}

func (p *part) GetBodyLen() int {
	return len(p.body)
}

func (p *part) GetBody() []byte {
	return common.CloneSlice(p.body)
}

func (p *part) WithBody(body []byte) Part {
	var exportHeaders headers.Headers
	if len(body) < 1 {
		exportHeaders = p.headers.WithoutHeader(headers.ContentType)
	} else {
		exportHeaders = p.headers.WithHeader(headers.ContentType, http.DetectContentType(body))
	}

	return &part{
		headers:  exportHeaders,
		body:     common.CloneSlice(body),
		subParts: p.subParts,
	}
}

func (p *part) WithBodyFromString(body string) Part {
	return p.WithBody([]byte(body))
}

func (p *part) GetSubParts() []Part {
	return p.subParts.ExtractPartsSlice()
}

func (p *part) WithSubParts(subParts ...Part) Part {
	return &part{
		headers:  p.headers,
		body:     p.body,
		subParts: NewPartsList(subParts...),
	}
}

func (p *part) Compile() ([]byte, error) {
	msg, err := p.ToPlainMessage()
	if err != nil {
		return nil, err
	}

	hds := headers.NewHeadersFromMap(msg.Header)
	bodyBytes, err := io.ReadAll(msg.Body)
	if err != nil {
		return nil, err
	}

	bytes := append(hds.Compile(), []byte("\r\n")...)
	bytes = append(bytes, bodyBytes...)

	return bytes, nil
}
