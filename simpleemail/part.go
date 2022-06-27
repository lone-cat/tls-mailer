package simpleemail

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	MultipartPrefix = `multipart/`

	MultipartMixed       = MultipartPrefix + `mixed`
	MultipartRelated     = MultipartPrefix + `related`
	MultipartAlternative = MultipartPrefix + `alternative`

	TextPlain = `text/plain`
	TextHtml  = `text/html`
)

type part struct {
	headers  *headers
	body     string
	subParts subParts
}

func newPart() *part {
	return &part{
		headers:  newHeaders(),
		subParts: newSubParts(),
	}
}

func newPartFromString(data string) (createdPart *part) {
	createdPart = newPart()
	createdPart = createdPart.withBody(data)
	return
}

func newEmbeddedPartFromString(cid, data string) (embeddedPart *part) {
	embeddedPart = newPartFromString(data)
	embeddedPart.headers = embeddedPart.headers.withHeader(ContentDispositionHeader, `inline`)
	if cid != `` {
		embeddedPart.headers = embeddedPart.headers.withHeader(ContentIdHeader, fmt.Sprintf(`<%s>`, cid))
	}

	return
}

func newEmbeddedPartFromFile(cid, filename string) (embedded *part, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	embedded = newEmbeddedPartFromString(cid, string(data))
	embedded.headers = embedded.headers.withHeader(ContentDispositionHeader, fmt.Sprintf(`inline; filename="%s"`, filepath.Base(filename)))

	return
}

func newAttachedPartFromString(data string) (attachment *part) {
	attachment = newPartFromString(data)
	attachment.headers = attachment.headers.withHeader(ContentDispositionHeader, `attachment`)

	return
}

func newAttachedPartFromFile(filename string) (attachment *part, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	attachment = newAttachedPartFromString(string(data))
	attachment.headers = attachment.headers.withHeader(ContentDispositionHeader, fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filename)))

	return
}

func (p *part) clone() *part {
	clonedPart := &part{
		headers:  p.headers,
		body:     p.body,
		subParts: p.subParts,
	}
	return clonedPart
}

func (p *part) getHeaders() *headers {
	return p.headers.clone()
}

func (p *part) withHeaders(headers *headers) *part {
	return &part{
		headers:  headers.clone(),
		body:     p.body,
		subParts: p.subParts,
	}
}

func (p *part) withHeadersFromMap(headers map[string][]string) *part {
	return &part{
		headers:  newHeadersFromMap(headers),
		body:     p.body,
		subParts: p.subParts,
	}
}

func (p *part) GetBody() string {
	return p.body
}

func (p *part) withBody(body string) *part {
	return &part{
		headers:  p.headers.withHeader(ContentTypeHeader, http.DetectContentType([]byte(body))),
		body:     body,
		subParts: p.subParts,
	}
}

func (p *part) getSubParts() subParts {
	return p.subParts.clone()
}

func (p *part) withSubParts(subPartsList subParts) *part {
	return &part{
		headers:  p.headers,
		body:     p.body,
		subParts: subPartsList.clone(),
	}
}

func (p *part) compile() ([]byte, error) {
	msg, err := p.toPlainMessage()
	if err != nil {
		return nil, err
	}

	headers := newHeadersFromMap(msg.Header)
	bodyBytes, err := io.ReadAll(msg.Body)
	if err != nil {
		return nil, err
	}

	bytes := append(headers.compile(), []byte("\r\n")...)
	bytes = append(bytes, bodyBytes...)

	return bytes, nil
}
