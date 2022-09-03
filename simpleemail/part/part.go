package part

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type part struct {
	headers  headers.Headers
	body     string
	subParts *partsList
}

func NewPart() Part {
	return &part{
		headers:  headers.NewHeaders(),
		subParts: newPartsList(),
	}
}

func NewPartFromString(data string) (createdPart Part) {
	return NewPart().WithBody(data)
}

func NewEmbeddedPartFromString(cid, data string) (embeddedPart Part) {
	embeddedPart = NewPartFromString(data)
	hdrs := embeddedPart.GetHeaders().
		WithHeader(
			headers.ContentDispositionHeader,
			`inline`,
		)
	if cid != `` {
		hdrs = hdrs.WithHeader(
			headers.ContentIdHeader,
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

	embedded = NewEmbeddedPartFromString(cid, string(data))
	hdrs := embedded.GetHeaders().
		WithHeader(
			headers.ContentDispositionHeader,
			fmt.Sprintf(`inline; filename="%s"`, filepath.Base(filename)),
		)

	embedded = embedded.WithHeaders(hdrs)

	return
}

func NewAttachedPartFromString(data string) (attachment Part) {
	attachment = NewPartFromString(data)
	hdrs := attachment.GetHeaders().
		WithHeader(
			headers.ContentDispositionHeader,
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

	attachment = NewAttachedPartFromString(string(data))
	hdrs := attachment.GetHeaders().
		WithHeader(
			headers.ContentDispositionHeader, fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filename)),
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

func (p *part) GetBody() string {
	return p.body
}

func (p *part) WithBody(body string) (exportPart Part) {
	var exportHeaders headers.Headers
	if body == `` {
		exportHeaders = p.headers.WithoutHeader(headers.ContentTypeHeader)
	} else {
		exportHeaders = p.headers.WithHeader(headers.ContentTypeHeader, http.DetectContentType([]byte(body)))
	}

	return &part{
		headers:  exportHeaders,
		body:     body,
		subParts: p.subParts,
	}
}

func (p *part) GetSubParts() *partsList {
	return p.subParts
}

func (p *part) WithSubParts(subPartsList *partsList) Part {
	return &part{
		headers:  p.headers,
		body:     p.body,
		subParts: subPartsList,
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
