package part

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail/encoding"
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
	if len(p.subParts.ExtractPartsSlice()) < 1 {

		var encodedBody []byte

		if p.GetBodyLen() > 0 {
			clonedHeaders = clonedHeaders.WithHeader(headers.ContentType, http.DetectContentType([]byte(p.GetBody())))

			var contentType string
			contentType, err = clonedHeaders.GetContentType()
			if err != nil {
				return
			}
			if strings.HasPrefix(contentType, `text/`) {
				//encodedBody = mime.QEncoding.Encode(`utf-8`, p.GetBody())
				encodedBody, err = encoding.BytesToQuotedPrintable(p.GetBody())
				if err != nil {
					return
				}
				clonedHeaders = clonedHeaders.WithHeader(headers.ContentTransferEncoding, encoding.QuotedPrintable.String())
			} else {
				//encodedBody = mime.BEncoding.Encode(`utf-8`, p.GetBody())
				encodedBody, err = encoding.BytesToBase64(p.GetBody())
				if err != nil {
					return nil, err
				}
				clonedHeaders = clonedHeaders.WithHeader(headers.ContentTransferEncoding, encoding.Base64.String())
			}
		}

		return &mail.Message{
			Header: clonedHeaders.ExtractHeadersMap(),
			Body:   bytes.NewReader(encodedBody),
		}, nil
	}

	subPartsSlice := p.subParts.ExtractPartsSlice()
	if len(subPartsSlice) == 1 {
		return subPartsSlice[0].ToPlainMessage()
	}

	contentType, params, _ := mime.ParseMediaType(clonedHeaders.GetFirstHeaderValue(headers.ContentType))
	needRegenerateHeader := false
	if !strings.HasPrefix(contentType, headers.MultipartPrefix) {
		contentType = headers.MultipartMixed
		needRegenerateHeader = true
	}
	boundary, exists := params[`boundary`]
	if !exists || boundary == `` {
		boundary = headers.GenerateBoundary()
		needRegenerateHeader = true
	}
	if needRegenerateHeader {
		clonedHeaders = clonedHeaders.WithHeader(headers.ContentType, fmt.Sprintf(`%s; boundary="%s"`, contentType, boundary))
		contentType, params, err = mime.ParseMediaType(clonedHeaders.GetFirstHeaderValue(headers.ContentType))
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

	for _, subPart := range p.subParts.ExtractPartsSlice() {
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
