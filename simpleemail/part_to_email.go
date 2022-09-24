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

func convertPartToEmail(sourcePart part.Part) (email *Email, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	email = NewEmptyEmail()

	email.mainPart, email.attachments, err = splitEmailPart(sourcePart)
	if err != nil {
		return
	}

	email.headers, email.from, email.to, email.cc, email.bcc, email.subject, err = proccessHeadersAndExtractPrimaryHeaders(sourcePart.GetHeaders())
	if err != nil {
		return
	}

	return
}

func splitEmailPart(prt part.Part) (relatedPart *relatedSubPart, attachments part.PartsList, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	var contentType string
	if prt.GetHeaders().GetFirstHeaderValue(headers.ContentTypeHeader) != `` || prt.GetBody() != `` {
		contentType, err = prt.GetHeaders().GetContentType()
		if err != nil {
			return
		}
	}

	attachments = part.NewPartsList()

	var partToConvert part.Part
	if contentType == headers.MultipartMixed {
		if len(prt.GetSubParts()) < 1 {
			relatedPart = newRelatedSubPart()
			return
		}
		partToConvert = prt.GetSubParts()[0]
		attachments = part.NewPartsList(prt.GetSubParts()[1:]...)
	} else {
		partToConvert = prt
	}

	relatedPart, err = convertToRelatedPart(partToConvert)
	if err != nil {
		return
	}

	return
}

func convertToRelatedPart(prt part.Part) (relatedPart *relatedSubPart, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	prtHeaders := prt.GetHeaders()
	var contentType string
	if prtHeaders.GetFirstHeaderValue(headers.ContentTypeHeader) != `` || prt.GetBody() != `` {
		contentType, err = prtHeaders.GetContentType()
		if err != nil {
			return
		}
	}

	relatedPart = newRelatedSubPart()

	var partToConvert part.Part
	if contentType == headers.MultipartRelated {
		prtSubParts := prt.GetSubParts()
		if len(prtSubParts) < 1 {
			return
		}
		partToConvert = prtSubParts[0]
		relatedPart.embeddedSubParts = part.NewPartsList(prtSubParts[1:]...)
		relatedPart.headers = prt.GetHeaders()
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

func convertToAlternativePart(prt part.Part) (alternativePart *alternativeSubPart, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	prtHeaders := prt.GetHeaders()
	var contentType string
	if prtHeaders.GetFirstHeaderValue(headers.ContentTypeHeader) != `` || prt.GetBody() != `` {
		contentType, err = prtHeaders.GetContentType()
		if err != nil {
			return
		}
		if prtHeaders.IsMultipart() && contentType != headers.MultipartAlternative {
			err = errors.New(fmt.Sprintf(`unexpected multipart type "%s"`, contentType))
			return
		}
	}

	alternativePart = newAlternativeSubPart()
	dataParts := []part.Part{prt}
	if contentType == headers.MultipartAlternative {
		prtSubParts := prt.GetSubParts()
		if len(prtSubParts) < 1 {
			return
		}
		if len(prtSubParts) > 2 {
			err = errors.New(`alternative part contains more than two subparts`)
			return
		}
		alternativePart.headers = prt.GetHeaders()
		dataParts = prt.GetSubParts()
	}

	var textPart, htmlPart part.Part
	var found bool

	textPart, found, err = extractOnePartByContentType(headers.TextPlain, dataParts...)
	if err != nil {
		return
	}
	if found {
		alternativePart.textPart = textPart
	}

	htmlPart, found, err = extractOnePartByContentType(headers.TextHtml, dataParts...)
	if err != nil {
		return
	}
	if found {
		alternativePart.htmlPart = htmlPart
	}

	return
}

func extractOnePartByContentType(contentType string, prts ...part.Part) (textPart part.Part, found bool, err error) {
	var subPartContentType string
	for _, prt := range prts {
		subPartContentType, err = prt.GetHeaders().GetContentType()
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
		textPart = prt.Clone()
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
