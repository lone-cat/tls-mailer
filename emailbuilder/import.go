package emailbuilder

import (
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

func Import(message string) (email Email, err error) {
	var msg *mail.Message
	msg, err = mail.ReadMessage(strings.NewReader(message))
	if err != nil {
		return
	}

	var part Part
	part, err = convertMessageToPartRecursive(msg)
	if err != nil {
		return
	}

	email.Text, email.Html, email.embedded, email.attachments, err = splitPart(part)
	if err != nil {
		return
	}

	email.headers, email.From, email.To, email.Cc, email.Bcc, email.Subject, err = proccessHeadersAndExtractPrimaryHeaders(part.GetHeaders())
	if err != nil {
		return
	}

	return
}

func convertMessageToPartRecursive(msg *mail.Message) (part Part, err error) {
	part = newPart()

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get(ContentTypeHeader))
	if err != nil {
		return
	}

	part = part.WithHeadersFromMap(msg.Header)

	if !strings.HasPrefix(mediaType, MultipartPrefix) {
		var msgBodyBytes []byte
		msgBodyBytes, err = io.ReadAll(msg.Body)
		if err != nil {
			return
		}
		part = part.WithBody(string(msgBodyBytes))

		transferEncoding := part.headers.GetContentTransferEncoding()
		if transferEncoding == EncodingQuotedPrintable || transferEncoding == EncodingBase64 {
			var partBody string
			partBody, err = extractBodyFromPart(part)
			if err != nil {
				return
			}
			part = part.WithBody(partBody)
			part = part.WithHeaders(part.GetHeaders().WithoutHeader(ContentTransferEncodingHeader))
		}

		return
	}

	subParts := NewSubParts()
	mr := multipart.NewReader(msg.Body, params["boundary"])
	var p *multipart.Part
	for {
		p, err = mr.NextRawPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		subMsg := &mail.Message{Header: mail.Header(p.Header), Body: p}
		var subPart Part
		subPart, err = convertMessageToPartRecursive(subMsg)
		if err != nil {
			return
		}
		subParts = append(subParts, subPart)
	}

	part = part.WithSubParts(subParts)

	return part, nil
}

func splitPart(part Part) (text string, html string, embedded []Part, attachments []Part, err error) {
	var contentType string
	contentType, err = part.GetHeaders().GetContentType()
	if err != nil {
		return
	}

	if contentType == MultipartMixed {
		text, html, embedded, attachments, err = splitMixedPart(part)
	}

	if contentType == MultipartRelated {
		text, html, embedded, err = splitRelatedPart(part)
		attachments = NewSubParts()
	}

	if contentType == MultipartAlternative {
		text, html, err = splitAlternativePart(part)
		embedded = NewSubParts()
		attachments = NewSubParts()
	}

	if contentType == TextPlain {
		text, err = extractBodyFromPart(part)
		embedded = NewSubParts()
		attachments = NewSubParts()
	}

	if contentType == TextHtml {
		html, err = extractBodyFromPart(part)
		embedded = NewSubParts()
		attachments = NewSubParts()

	}

	return
}

func splitMixedPart(part Part) (text string, html string, embedded []Part, attachments []Part, err error) {
	parts := part.GetSubParts()

	if len(parts) < 1 {
		return ``, ``, nil, nil, errors.New(`mixed part is empty`)
	}

	var contentType string
	contentType, err = parts[0].GetHeaders().GetContentType()
	if err != nil {
		return
	}

	if contentType == MultipartRelated {
		text, html, embedded, err = splitRelatedPart(parts[0])
		if err != nil {
			return
		}
		attachments = parts[1:]
		return
	}

	if contentType == MultipartAlternative {
		text, html, err = splitAlternativePart(parts[0])
		if err != nil {
			return
		}
		embedded = NewSubParts()
		attachments = parts[1:]
		return
	}

	if contentType == TextPlain {
		text, err = extractBodyFromPart(parts[0])
		if err != nil {
			return
		}
		embedded = NewSubParts()
		parts = parts[1:]
		attachments = parts
	}

	if len(parts) < 1 {
		return
	}

	if contentType == TextHtml {
		html, err = extractBodyFromPart(parts[0])
		if err != nil {
			return
		}
		embedded = NewSubParts()
		parts = parts[1:]
		attachments = parts
	}

	return
}

func splitRelatedPart(part Part) (text string, html string, embedded []Part, err error) {
	parts := part.GetSubParts()

	if len(parts) < 1 {
		return ``, ``, nil, errors.New(`related part is empty`)
	}

	var contentType string
	contentType, err = parts[0].GetHeaders().GetContentType()
	if err != nil {
		return
	}

	if contentType == MultipartAlternative {
		text, html, err = splitAlternativePart(parts[0])
		if err != nil {
			return
		}
		embedded = parts[1:]
		return
	}

	if contentType == TextPlain {
		text, err = extractBodyFromPart(parts[0])
		if err != nil {
			return
		}
		parts = parts[1:]
	}

	if len(parts) < 1 {
		embedded = NewSubParts()
		return
	}

	if contentType == TextHtml {
		html, err = extractBodyFromPart(parts[0])
		if err != nil {
			return
		}
		parts = parts[1:]
	}

	embedded = parts

	return
}

func splitAlternativePart(part Part) (text string, html string, err error) {
	parts := part.GetSubParts()

	if len(parts) < 1 {
		return ``, ``, errors.New(`alternative part is empty`)

	}
	if len(parts) > 2 {
		return ``, ``, errors.New(`alternative part includes more than 2 parts`)
	}

	attachments := NewSubParts()
	for _, part := range parts {
		var contentType string
		contentType, err = part.GetHeaders().GetContentType()
		if err != nil {
			return
		}

		if contentType == TextPlain {
			if text == `` {
				text, err = extractBodyFromPart(part)
				if err != nil {
					return ``, "", err
				}
				continue
			} else {
				return ``, ``, errors.New(`more than one text part in alternative block`)
			}
		} else if contentType == TextHtml {
			if html == `` {
				html, err = extractBodyFromPart(part)
				if err != nil {
					return ``, ``, err
				}
				continue
			} else {
				return ``, ``, errors.New(`more than one text part in alternative block`)
			}
		} else {
			return ``, ``, errors.New(`unsupported alternative part "` + contentType + `"`)
		}
	}

	if len(attachments) > 0 {
		return ``, ``, errors.New(`alternative part includes unexpected part types`)
	}

	return
}

func extractBodyFromPart(part Part) (decodedBody string, err error) {
	encoding := part.GetHeaders().GetContentTransferEncoding()
	if encoding == EncodingEmpty || encoding == Encoding7bit || encoding == Encoding8bit || encoding == EncodingBinary {
		return part.GetBody(), nil
	}

	if encoding == EncodingQuotedPrintable {
		var decodedBodyBytes []byte
		decodedBodyBytes, err = io.ReadAll(quotedprintable.NewReader(strings.NewReader(part.GetBody())))
		if err != nil {
			return
		}
		decodedBody = string(decodedBodyBytes)
		return
	}
	if encoding == EncodingBase64 {
		var decodedBodyBytes []byte
		decodedBodyBytes, err = io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(part.GetBody())))
		if err != nil {
			return
		}
		decodedBody = string(decodedBodyBytes)
		return
	}
	return ``, errors.New(`unexpected body encoding`)
}

func proccessHeadersAndExtractPrimaryHeaders(oldHeaders Headers) (headers Headers, from []mail.Address, to []mail.Address, cc []mail.Address, bcc []mail.Address, subject string, err error) {
	headers = oldHeaders.Clone()

	from = make([]mail.Address, 0)
	to = make([]mail.Address, 0)
	cc = make([]mail.Address, 0)
	bcc = make([]mail.Address, 0)

	if headers.GetFirstHeaderValue(FromHeader) != `` {
		from, err = headers.GetAddressList(FromHeader)
		if err != nil {
			return
		}
	}

	if headers.GetFirstHeaderValue(ToHeader) != `` {
		to, err = headers.GetAddressList(ToHeader)
		if err != nil {
			return
		}
	}

	if headers.GetFirstHeaderValue(CCHeader) != `` {
		cc, err = headers.GetAddressList(CCHeader)
		if err != nil {
			return
		}
	}

	if headers.GetFirstHeaderValue(BCCHeader) != `` {
		bcc, err = headers.GetAddressList(BCCHeader)
		if err != nil {
			return
		}
	}

	subject = headers.GetFirstHeaderValue(SubjectHeader)
	if subject != `` {
		subject, err = decoder.DecodeHeader(subject)
		if err != nil {
			return
		}
	}

	headers = headers.WithoutHeader(FromHeader).
		WithoutHeader(ToHeader).
		WithoutHeader(CCHeader).
		WithoutHeader(BCCHeader).
		WithoutHeader(SubjectHeader)

	return
}
