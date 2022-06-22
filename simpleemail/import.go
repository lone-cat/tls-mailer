package simpleemail

import (
	"encoding/base64"
	"errors"
	"github.com/lone-cat/stackerrors"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

func Import(message string) (email Email, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(`Import`, err)
		return
	}()

	var msg *mail.Message
	msg, err = mail.ReadMessage(strings.NewReader(message))
	if err != nil {
		return
	}

	var convertedPart part
	convertedPart, err = convertMessageToPartRecursive(msg)
	if err != nil {
		return
	}

	email.mainPart, email.attachments, err = splitPart(convertedPart)
	if err != nil {
		return
	}

	email.headers, email.from, email.to, email.cc, email.bcc, email.subject, err = proccessHeadersAndExtractPrimaryHeaders(convertedPart.getHeaders())
	if err != nil {
		return
	}

	return
}

func convertMessageToPartRecursive(msg *mail.Message) (exportedPart part, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(`convertMessageToPartRecursive`, err)
		return
	}()
	exportedPart = newPart()

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get(ContentTypeHeader))
	if err != nil {
		return
	}

	exportedPart = exportedPart.withHeadersFromMap(msg.Header)

	if !strings.HasPrefix(mediaType, MultipartPrefix) {
		var msgBodyBytes []byte
		msgBodyBytes, err = io.ReadAll(msg.Body)
		if err != nil {
			return
		}
		exportedPart = exportedPart.withBody(string(msgBodyBytes))

		transferEncoding := exportedPart.headers.getContentTransferEncoding()
		if transferEncoding == EncodingQuotedPrintable || transferEncoding == EncodingBase64 {
			var partBody string
			partBody, err = extractBodyFromPart(exportedPart)
			if err != nil {
				return
			}
			exportedPart = exportedPart.withBody(partBody)
			exportedPart = exportedPart.withHeaders(exportedPart.getHeaders().withoutHeader(ContentTransferEncodingHeader))
		}

		return
	}

	convertedSubParts := newSubParts()
	mr := multipart.NewReader(msg.Body, params["boundary"])
	var p *multipart.Part
	for {
		p, err = mr.NextRawPart()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		subMsg := &mail.Message{Header: mail.Header(p.Header), Body: p}
		var subPart part
		subPart, err = convertMessageToPartRecursive(subMsg)
		if err != nil {
			return
		}
		convertedSubParts = append(convertedSubParts, subPart)
	}

	exportedPart = exportedPart.withSubParts(convertedSubParts)

	return
}

func splitPart(part part) (mainPart mainSubPart, attachments subParts, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(`splitPart`, err)
		return
	}()
	var contentType string
	contentType, err = part.getHeaders().getContentType()
	if err != nil {
		return
	}

	var text, html string
	embedded := newSubParts()
	attachments = newSubParts()

	if contentType == MultipartMixed {
		text, html, embedded, attachments, err = splitMixedPart(part)
	}

	if contentType == MultipartRelated {
		text, html, embedded, err = splitRelatedPart(part)
	}

	if contentType == MultipartAlternative {
		text, html, err = splitAlternativePart(part)
	}

	if contentType == TextPlain {
		text, err = extractBodyFromPart(part)
	}

	if contentType == TextHtml {
		html, err = extractBodyFromPart(part)
	}

	mainPart.textSubPart = mainPart.textSubPart.withText(text)
	mainPart.textSubPart = mainPart.textSubPart.withHtml(html)
	mainPart.embeddedSubParts = embedded

	return
}

func splitMixedPart(part part) (text string, html string, embedded subParts, attachments subParts, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(`splitMixedPart`, err)
		return
	}()
	parts := part.getSubParts()

	if len(parts) < 1 {
		return ``, ``, nil, nil, errors.New(`mixed part is empty`)
	}

	var contentType string
	contentType, err = parts[0].getHeaders().getContentType()
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
		embedded = newSubParts()
		attachments = parts[1:]
		return
	}

	if contentType == TextPlain {
		text, err = extractBodyFromPart(parts[0])
		if err != nil {
			return
		}
		embedded = newSubParts()
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
		embedded = newSubParts()
		parts = parts[1:]
		attachments = parts
	}

	return
}

func splitRelatedPart(part part) (text string, html string, embedded subParts, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(`splitRelatedPart`, err)
		return
	}()
	parts := part.getSubParts()

	if len(parts) < 1 {
		return ``, ``, nil, errors.New(`related part is empty`)
	}

	var contentType string
	contentType, err = parts[0].getHeaders().getContentType()
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
		embedded = newSubParts()
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

func splitAlternativePart(part part) (text string, html string, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(`splitAlternativePart`, err)
		return
	}()
	parts := part.getSubParts()

	if len(parts) < 1 {
		return ``, ``, errors.New(`alternative part is empty`)
	}
	if len(parts) > 2 {
		return ``, ``, errors.New(`alternative part includes more than 2 parts`)
	}

	attachments := newSubParts()
	for _, part := range parts {
		var contentType string
		contentType, err = part.getHeaders().getContentType()
		if err != nil {
			return
		}

		if contentType == TextPlain {
			if text == `` {
				text, err = extractBodyFromPart(part)
				if err != nil {
					return ``, ``, err
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

func extractBodyFromPart(part part) (decodedBody string, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(`extractBodyFromPart`, err)
		return
	}()
	encoding := part.getHeaders().getContentTransferEncoding()
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
	defer func() {
		err = stackerrors.WrapInDefer(`proccessHeadersAndExtractPrimaryHeaders`, err)
		return
	}()
	headers = oldHeaders.clone()

	from = make([]mail.Address, 0)
	to = make([]mail.Address, 0)
	cc = make([]mail.Address, 0)
	bcc = make([]mail.Address, 0)

	if headers.getFirstHeaderValue(FromHeader) != `` {
		from, err = headers.getAddressList(FromHeader)
		if err != nil {
			return
		}
	}

	if headers.getFirstHeaderValue(ToHeader) != `` {
		to, err = headers.getAddressList(ToHeader)
		if err != nil {
			return
		}
	}

	if headers.getFirstHeaderValue(CCHeader) != `` {
		cc, err = headers.getAddressList(CCHeader)
		if err != nil {
			return
		}
	}

	if headers.getFirstHeaderValue(BCCHeader) != `` {
		bcc, err = headers.getAddressList(BCCHeader)
		if err != nil {
			return
		}
	}

	subject = headers.getFirstHeaderValue(SubjectHeader)
	if subject != `` {
		subject, err = decoder.DecodeHeader(subject)
		if err != nil {
			return
		}
	}

	headers = headers.withoutHeader(FromHeader).
		withoutHeader(ToHeader).
		withoutHeader(CCHeader).
		withoutHeader(BCCHeader).
		withoutHeader(SubjectHeader)

	return
}
