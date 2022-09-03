package part

import (
	"errors"
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail/encode"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/textproto"
	"strings"
)

func (p *part) ToPlainMessage() (msg *mail.Message, err error) {
	clonedHeaders := p.GetHeaders()
	if len(p.subParts.parts) < 1 {

		var encodedBody string

		if p.GetBody() != `` {
			clonedHeaders = clonedHeaders.WithHeader(headers.ContentTypeHeader, http.DetectContentType([]byte(p.GetBody())))

			var contentType string
			contentType, err = clonedHeaders.GetContentType()
			if err != nil {
				return
			}
			if strings.HasPrefix(contentType, `text/`) {
				//encodedBody = mime.QEncoding.Encode(`utf-8`, p.GetBody())
				encodedBody, err = encode.ToQuotedPrintable(p.GetBody())
				if err != nil {
					return
				}
				clonedHeaders = clonedHeaders.WithHeader(headers.ContentTransferEncodingHeader, headers.EncodingQuotedPrintable.String())
			} else {
				//encodedBody = mime.BEncoding.Encode(`utf-8`, p.GetBody())
				encodedBody, err = encode.ToBase64(p.GetBody())
				if err != nil {
					return nil, err
				}
				clonedHeaders = clonedHeaders.WithHeader(headers.ContentTransferEncodingHeader, headers.EncodingBase64.String())
			}
		}

		return &mail.Message{
			Header: clonedHeaders.ExtractHeadersMap(),
			Body:   strings.NewReader(encodedBody),
		}, nil
	}

	if len(p.subParts.parts) == 1 {
		return p.subParts.parts[0].ToPlainMessage()
	}

	contentType, params, _ := mime.ParseMediaType(clonedHeaders.GetFirstHeaderValue(headers.ContentTypeHeader))
	needRegenerateHeader := false
	if !strings.HasPrefix(contentType, MultipartPrefix) {
		contentType = MultipartMixed
		needRegenerateHeader = true
	}
	boundary, exists := params[`boundary`]
	if !exists || boundary == `` {
		boundary = headers.GenerateBoundary()
		needRegenerateHeader = true
	}
	if needRegenerateHeader {
		clonedHeaders = clonedHeaders.WithHeader(headers.ContentTypeHeader, fmt.Sprintf(`%s; boundary="%s"`, contentType, boundary))
		contentType, params, err = mime.ParseMediaType(clonedHeaders.GetFirstHeaderValue(headers.ContentTypeHeader))
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

	for _, subPart := range p.subParts.parts {
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
