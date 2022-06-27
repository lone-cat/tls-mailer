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
	"strings"
)

func (p *part) toPlainMessage() (msg *mail.Message, err error) {
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
