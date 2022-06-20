package simpleemail

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
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

func newEmbeddedPart(cid, filename string) (embedded part, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	headers := newHeaders()
	headers = headers.withHeader(`content-disposition`, fmt.Sprintf(`inline; filename="%s"`, filepath.Base(filename)))
	if cid != `` {
		headers = headers.withHeader(`content-id`, fmt.Sprintf(`<%s>`, cid))
	}

	embedded = part{
		headers: headers,
		body:    string(data),
	}

	return
}

func newAttachedPart(filename string) (embedded part, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	headers := newHeaders()
	headers = headers.withHeader(`content-disposition`, fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filename)))

	embedded = part{
		headers: headers,
		body:    string(data),
	}

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

func (p part) toPlainMessage() (msg *mail.Message, err error) {
	clonedHeaders := p.getHeaders()
	if len(p.subParts) < 1 {

		var encodedBody string

		if p.GetBody() != `` {
			clonedHeaders = clonedHeaders.withHeader(ContentTypeHeader, http.DetectContentType([]byte(p.GetBody())))

			var contentType string
			contentType, err = clonedHeaders.getContentType()
			if err != nil {
				return
			}
			if strings.HasPrefix(contentType, `text/`) {
				//encodedBody = mime.QEncoding.Encode(`utf-8`, p.GetBody())
				encodedBody, err = toQuotedPrintable(p.GetBody())
				if err != nil {
					return
				}
				clonedHeaders = clonedHeaders.withHeader(ContentTransferEncodingHeader, EncodingQuotedPrintable.String())
			} else {
				//encodedBody = mime.BEncoding.Encode(`utf-8`, p.GetBody())
				encodedBody, err = toBase64(p.GetBody())
				if err != nil {
					return nil, err
				}
				clonedHeaders = clonedHeaders.withHeader(ContentTransferEncodingHeader, EncodingBase64.String())
			}
		}

		return &mail.Message{
			Header: clonedHeaders.extractHeadersMap(),
			Body:   strings.NewReader(encodedBody),
		}, nil
	}

	if len(p.subParts) == 1 {
		return p.subParts[0].toPlainMessage()
	}

	contentType, params, _ := mime.ParseMediaType(clonedHeaders.getFirstHeaderValue(ContentTypeHeader))
	needRegenerateHeader := false
	if !strings.HasPrefix(contentType, MultipartPrefix) {
		contentType = MultipartMixed
		needRegenerateHeader = true
	}
	boundary, exists := params[`boundary`]
	if !exists || boundary == `` {
		boundary = generateBoundary()
		needRegenerateHeader = true
	}
	if needRegenerateHeader {
		clonedHeaders = clonedHeaders.withHeader(ContentTypeHeader, fmt.Sprintf(`%s; boundary="%s"`, contentType, boundary))
		contentType, params, err = mime.ParseMediaType(clonedHeaders.getFirstHeaderValue(ContentTypeHeader))
		if err != nil {
			return
		}
		boundary, exists = params[`boundary`]
		if !exists || boundary == `` {
			return nil, errors.New(`boundary not set`)
		}
	}

	bodyWriter := &strings.Builder{}
	multipartWriter := multipart.NewWriter(bodyWriter)
	err = multipartWriter.SetBoundary(boundary)
	if err != nil {
		return
	}

	var subMsg *mail.Message
	var partWriter io.Writer
	var subMsgBodyBytes []byte

	for _, subPart := range p.subParts {
		subMsg, err = subPart.toPlainMessage()
		if err != nil {
			return
		}
		partWriter, err = multipartWriter.CreatePart(textproto.MIMEHeader(subMsg.Header))
		if err != nil {
			return
		}
		subMsgBodyBytes, err = io.ReadAll(subMsg.Body)
		if err != nil {
			return
		}
		_, err = partWriter.Write(subMsgBodyBytes)
		if err != nil {
			return
		}
	}
	err = multipartWriter.Close()
	if err != nil {
		return
	}

	return &mail.Message{
		Header: clonedHeaders.extractHeadersMap(),
		Body:   strings.NewReader(bodyWriter.String()),
	}, nil
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
