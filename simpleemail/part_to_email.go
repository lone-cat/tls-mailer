package simpleemail

import (
	"errors"
	"fmt"
	"github.com/lone-cat/stackerrors"
	"net/mail"
)

func convertPartToEmail(sourcePart *part) (email *Email, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	email = NewEmptyEmail()

	email.mainPart, email.attachments, err = splitEmailPart(sourcePart)
	if err != nil {
		return
	}

	email.headers, email.from, email.to, email.cc, email.bcc, email.subject, err = proccessHeadersAndExtractPrimaryHeaders(sourcePart.getHeaders())
	if err != nil {
		return
	}

	return
}

func splitEmailPart(prt *part) (relatedPart *relatedSubPart, attachments subParts, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	var contentType string
	if prt.headers.getFirstHeaderValue(ContentTypeHeader) != `` || prt.body != `` {
		contentType, err = prt.headers.getContentType()
		if err != nil {
			return
		}
	}

	attachments = newSubParts()

	var partToConvert *part
	if contentType == MultipartMixed {
		if len(prt.subParts) < 1 {
			relatedPart = newRelatedSubPart()
			return
		}
		partToConvert = prt.subParts[0]
		attachments = prt.subParts[1:].clone()
	} else {
		partToConvert = prt
	}

	relatedPart, err = convertToRelatedPart(partToConvert)
	if err != nil {
		return
	}

	return
}

func convertToRelatedPart(prt *part) (relatedPart *relatedSubPart, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	var contentType string
	if prt.headers.getFirstHeaderValue(ContentTypeHeader) != `` || prt.body != `` {
		contentType, err = prt.headers.getContentType()
		if err != nil {
			return
		}
	}

	relatedPart = newRelatedSubPart()

	var partToConvert *part
	if contentType == MultipartRelated {
		if len(prt.subParts) < 1 {
			return
		}
		partToConvert = prt.subParts[0]
		relatedPart.embeddedSubParts = prt.subParts[1:].clone()
		relatedPart.headers = prt.headers
	} else {
		partToConvert = prt
	}

	altPart, err := convertToAlternativePart(partToConvert)
	if err != nil {
		return
	}

	relatedPart.alternativeSubPart = altPart

	return
}

func convertToAlternativePart(prt *part) (alternativePart *alternativeSubPart, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	var contentType string
	if prt.headers.getFirstHeaderValue(ContentTypeHeader) != `` || prt.body != `` {
		contentType, err = prt.headers.getContentType()
		if err != nil {
			return
		}
		if prt.headers.isMultipart() && contentType != MultipartAlternative {
			err = errors.New(fmt.Sprintf(`unexpected multipart type "%s"`, contentType))
			return
		}
	}

	alternativePart = newAlternativeSubPart()
	dataParts := []*part{prt}
	if contentType == MultipartAlternative {
		if len(prt.subParts) < 1 {
			return
		}
		if len(prt.subParts) > 2 {
			err = errors.New(`alternative part contains more than two subparts`)
			return
		}
		alternativePart.headers = prt.headers
		dataParts = prt.subParts
	}

	var textPart, htmlPart *part
	var found bool

	textPart, found, err = extractOnePartByContentType(TextPlain, dataParts...)
	if err != nil {
		return
	}
	if found {
		alternativePart.textPart = textPart
	}

	htmlPart, found, err = extractOnePartByContentType(TextHtml, dataParts...)
	if err != nil {
		return
	}
	if found {
		alternativePart.htmlPart = htmlPart
	}

	return
}

func extractOnePartByContentType(contentType string, prts ...*part) (textPart *part, found bool, err error) {
	var subPartContentType string
	for _, prt := range prts {
		subPartContentType, err = prt.getHeaders().getContentType()
		if err != nil {
			err = nil
			continue
		}
		if subPartContentType != contentType {
			continue
		}
		if found {
			err = errors.New(`two parts with same content type found`)
			return
		}
		textPart = prt.clone()
		found = true
	}

	return
}

func proccessHeadersAndExtractPrimaryHeaders(oldHeaders *headers) (headers *headers, from []*mail.Address, to []*mail.Address, cc []*mail.Address, bcc []*mail.Address, subject string, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	from = make([]*mail.Address, 0)
	to = make([]*mail.Address, 0)
	cc = make([]*mail.Address, 0)
	bcc = make([]*mail.Address, 0)

	if oldHeaders.getFirstHeaderValue(FromHeader) != `` {
		from, err = oldHeaders.headers.AddressList(FromHeader)
		if err != nil {
			return
		}
	}

	if oldHeaders.getFirstHeaderValue(ToHeader) != `` {
		to, err = oldHeaders.headers.AddressList(ToHeader)
		if err != nil {
			return
		}
	}

	if oldHeaders.getFirstHeaderValue(CCHeader) != `` {
		cc, err = oldHeaders.headers.AddressList(CCHeader)
		if err != nil {
			return
		}
	}

	if oldHeaders.getFirstHeaderValue(BCCHeader) != `` {
		bcc, err = oldHeaders.headers.AddressList(BCCHeader)
		if err != nil {
			return
		}
	}

	subject = oldHeaders.getFirstHeaderValue(SubjectHeader)
	if subject != `` {
		subject, err = decoder.DecodeHeader(subject)
		if err != nil {
			return
		}
	}

	headers = oldHeaders.withoutHeader(FromHeader).
		withoutHeader(ToHeader).
		withoutHeader(CCHeader).
		withoutHeader(BCCHeader).
		withoutHeader(SubjectHeader)

	return
}
