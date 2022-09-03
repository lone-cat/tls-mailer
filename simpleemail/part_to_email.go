package simpleemail

import (
	"errors"
	"fmt"
	"github.com/lone-cat/stackerrors"
	"github.com/lone-cat/tls-mailer/simpleemail/encode"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"github.com/lone-cat/tls-mailer/simpleemail/part"
	"net/mail"
)

func convertPartToEmail(sourcePart *part.part) (email *Email, err error) {
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

func splitEmailPart(prt *part.part) (relatedPart *relatedSubPart, attachments part.subParts, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	var contentType string
	if prt.headers.GetFirstHeaderValue(headers.ContentTypeHeader) != `` || prt.body != `` {
		contentType, err = prt.headers.GetContentType()
		if err != nil {
			return
		}
	}

	attachments = part.newSubParts()

	var partToConvert *part.part
	if contentType == part.MultipartMixed {
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

func convertToRelatedPart(prt *part.part) (relatedPart *relatedSubPart, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	var contentType string
	if prt.headers.GetFirstHeaderValue(headers.ContentTypeHeader) != `` || prt.body != `` {
		contentType, err = prt.headers.GetContentType()
		if err != nil {
			return
		}
	}

	relatedPart = newRelatedSubPart()

	var partToConvert *part.part
	if contentType == part.MultipartRelated {
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

func convertToAlternativePart(prt *part.part) (alternativePart *alternativeSubPart, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	var contentType string
	if prt.headers.GetFirstHeaderValue(headers.ContentTypeHeader) != `` || prt.body != `` {
		contentType, err = prt.headers.GetContentType()
		if err != nil {
			return
		}
		if prt.headers.IsMultipart() && contentType != part.MultipartAlternative {
			err = errors.New(fmt.Sprintf(`unexpected multipart type "%s"`, contentType))
			return
		}
	}

	alternativePart = newAlternativeSubPart()
	dataParts := []*part.part{prt}
	if contentType == part.MultipartAlternative {
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

	var textPart, htmlPart *part.part
	var found bool

	textPart, found, err = extractOnePartByContentType(part.TextPlain, dataParts...)
	if err != nil {
		return
	}
	if found {
		alternativePart.textPart = textPart
	}

	htmlPart, found, err = extractOnePartByContentType(part.TextHtml, dataParts...)
	if err != nil {
		return
	}
	if found {
		alternativePart.htmlPart = htmlPart
	}

	return
}

func extractOnePartByContentType(contentType string, prts ...*part.part) (textPart *part.part, found bool, err error) {
	var subPartContentType string
	for _, prt := range prts {
		subPartContentType, err = prt.headers.GetContentType()
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

func proccessHeadersAndExtractPrimaryHeaders(oldHeaders headers.Headers) (hds headers.Headers, from []*mail.Address, to []*mail.Address, cc []*mail.Address, bcc []*mail.Address, subject string, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	from = make([]*mail.Address, 0)
	to = make([]*mail.Address, 0)
	cc = make([]*mail.Address, 0)
	bcc = make([]*mail.Address, 0)

	if oldHeaders.GetFirstHeaderValue(headers.FromHeader) != `` {
		from, err = oldHeaders.AddressList(headers.FromHeader)
		if err != nil {
			return
		}
	}

	if oldHeaders.GetFirstHeaderValue(headers.ToHeader) != `` {
		to, err = oldHeaders.AddressList(headers.ToHeader)
		if err != nil {
			return
		}
	}

	if oldHeaders.GetFirstHeaderValue(headers.CCHeader) != `` {
		cc, err = oldHeaders.AddressList(headers.CCHeader)
		if err != nil {
			return
		}
	}

	if oldHeaders.GetFirstHeaderValue(headers.BCCHeader) != `` {
		bcc, err = oldHeaders.AddressList(headers.BCCHeader)
		if err != nil {
			return
		}
	}

	subject = oldHeaders.GetFirstHeaderValue(headers.SubjectHeader)
	if subject != `` {
		subject, err = encode.DecodeHeader(subject)
		if err != nil {
			return
		}
	}

	hds = oldHeaders.WithoutHeader(headers.FromHeader).
		WithoutHeader(headers.ToHeader).
		WithoutHeader(headers.CCHeader).
		WithoutHeader(headers.BCCHeader).
		WithoutHeader(headers.SubjectHeader)

	return
}
