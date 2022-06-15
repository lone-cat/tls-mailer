package emailbuilder

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/textproto"
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

type Part struct {
	headers  Headers
	body     string
	subParts subParts
}

func newPart() Part {
	return Part{
		headers:  NewHeaders(),
		subParts: NewSubParts(),
	}
}

func (p Part) Clone() Part {
	p.headers = p.headers.Clone()
	p.subParts = p.subParts.Clone()
	return p
}

func (p Part) GetHeaders() Headers {
	return p.headers.Clone()
}

func (p Part) WithHeaders(headers Headers) Part {
	part := p.Clone()
	part.headers = headers.Clone()
	return part
}

func (p Part) WithHeadersFromMap(headers map[string][]string) Part {
	part := p.Clone()
	part.headers = NewHeadersFromMap(headers)
	return part
}

func (p Part) GetBody() string {
	return p.body
}

func (p Part) WithBody(body string) Part {
	part := p.Clone()
	part.body = body
	return part
}

func (p Part) GetSubParts() []Part {
	return p.subParts.Clone()
}

func (p Part) WithSubParts(subPartsList []Part) Part {
	part := p.Clone()
	part.subParts = subParts(subPartsList).Clone()
	return part
}

func (p Part) ToPlainMessage() (msg *mail.Message, err error) {
	clonedHeaders := p.GetHeaders()
	if len(p.subParts) < 1 {
		clonedHeaders = clonedHeaders.WithHeader(ContentTypeHeader, http.DetectContentType([]byte(p.GetBody())))

		var contentType string
		contentType, err = clonedHeaders.GetContentType()
		if err != nil {
			return
		}
		var encodedBody string
		if strings.HasPrefix(contentType, `text/`) {
			//encodedBody = mime.QEncoding.Encode(`utf-8`, p.GetBody())
			encodedBody, err = toQuotedPrintable(p.GetBody())
			if err != nil {
				return
			}
			clonedHeaders = clonedHeaders.WithHeader(ContentTransferEncodingHeader, EncodingQuotedPrintable.String())
		} else {
			//encodedBody = mime.BEncoding.Encode(`utf-8`, p.GetBody())
			encodedBody, err = toBase64(p.GetBody())
			if err != nil {
				return nil, err
			}
			clonedHeaders = clonedHeaders.WithHeader(ContentTransferEncodingHeader, EncodingBase64.String())
		}

		return &mail.Message{
			Header: clonedHeaders.ExtractHeadersMap(),
			Body:   strings.NewReader(encodedBody),
		}, nil
	}

	if len(p.subParts) == 1 {
		return p.subParts[0].ToPlainMessage()
	}

	contentType, params, _ := mime.ParseMediaType(clonedHeaders.GetFirstHeaderValue(ContentTypeHeader))
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
		clonedHeaders = clonedHeaders.WithHeader(ContentTypeHeader, fmt.Sprintf(`%s; boundary="%s"`, contentType, boundary))
		contentType, params, err = mime.ParseMediaType(clonedHeaders.GetFirstHeaderValue(ContentTypeHeader))
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
		subMsg, err = subPart.ToPlainMessage()
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
		Header: clonedHeaders.ExtractHeadersMap(),
		Body:   strings.NewReader(bodyWriter.String()),
	}, nil
}

func (p Part) Render() (string, error) {
	msg, err := p.ToPlainMessage()
	if err != nil {
		return ``, err
	}

	headers := NewHeadersFromMap(msg.Header)
	bodyBytes, err := io.ReadAll(msg.Body)
	if err != nil {
		return ``, err
	}

	return headers.Render() + "\r\n" + string(bodyBytes), nil
}
