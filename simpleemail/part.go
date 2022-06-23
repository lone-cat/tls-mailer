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
	headers  Headers
	body     string
	subParts subParts
}

func newPart() part {
	return part{
		headers:  newHeaders(),
		subParts: newSubParts(),
	}
}

func newPartFromString(data string) (createdPart part) {
	createdPart = newPart()
	createdPart = createdPart.withBody(data)
	return
}

func newEmbeddedPartFromString(cid, data string) (embeddedPart part) {
	embeddedPart = newPartFromString(data)
	embeddedPart.headers = embeddedPart.headers.withHeader(ContentDispositionHeader, `inline`)
	if cid != `` {
		embeddedPart.headers = embeddedPart.headers.withHeader(ContentIdHeader, fmt.Sprintf(`<%s>`, cid))
	}

	return
}

func newEmbeddedPartFromFile(cid, filename string) (embedded part, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	embedded = newEmbeddedPartFromString(cid, string(data))
	embedded.headers = embedded.headers.withHeader(ContentDispositionHeader, fmt.Sprintf(`inline; filename="%s"`, filepath.Base(filename)))

	return
}

func newAttachedPartFromString(data string) (attachment part) {
	attachment = newPartFromString(data)
	attachment.headers = attachment.headers.withHeader(ContentDispositionHeader, `attachment`)

	return
}

func newAttachedPartFromFile(filename string) (attachment part, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	attachment = newAttachedPartFromString(string(data))
	attachment.headers = attachment.headers.withHeader(ContentDispositionHeader, fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filename)))

	return
}

func (p part) clone() part {
	clonedPart := newPart()
	clonedPart.headers = p.headers.clone()
	clonedPart.body = p.body
	clonedPart.subParts = p.subParts.clone()
	return clonedPart
}

func (p part) getHeaders() Headers {
	return p.headers.clone()
}

func (p part) withHeaders(headers Headers) part {
	clonedPart := p.clone()
	clonedPart.headers = headers.clone()
	return clonedPart
}

func (p part) withHeadersFromMap(headers map[string][]string) part {
	clonedPart := p.clone()
	clonedPart.headers = newHeadersFromMap(headers)
	return clonedPart
}

func (p part) GetBody() string {
	return p.body
}

func (p part) withBody(body string) part {
	clonedPart := p.clone()
	clonedPart.body = body
	clonedPart.headers = clonedPart.headers.withHeader(ContentTypeHeader, http.DetectContentType([]byte(body)))
	return clonedPart
}

func (p part) getSubParts() subParts {
	return p.subParts.clone()
}

func (p part) withSubParts(subPartsList subParts) part {
	clonedPart := p.clone()
	clonedPart.subParts = subPartsList.clone()
	return clonedPart
}

func (p part) compile() ([]byte, error) {
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
