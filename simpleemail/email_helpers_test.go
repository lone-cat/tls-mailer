package simpleemail

import (
	"fmt"
	"github.com/lone-cat/stackerrors"
	"net/mail"
	"strings"
)

func emailsDiffErrors(e1 *Email, e2 *Email) []error {
	errors := make([]error, 0)
	if !addressSlicesEqual(e1.from, e2.from) {
		errors = append(
			errors,
			stackerrors.Newf(
				"`From` does not match: expected %s, got %s",
				convertAddressListToReadable(e1.from),
				convertAddressListToReadable(e2.from),
			),
		)
	}
	if !addressSlicesEqual(e1.to, e2.to) {
		errors = append(
			errors,
			stackerrors.Newf(
				"`To` does not match: expected %s, got %s",
				convertAddressListToReadable(e1.to),
				convertAddressListToReadable(e2.to),
			),
		)
	}
	if !addressSlicesEqual(e1.cc, e2.cc) {
		errors = append(
			errors,
			stackerrors.Newf(
				"`Cc` does not match: expected %s, got %s",
				convertAddressListToReadable(e1.cc),
				convertAddressListToReadable(e2.cc),
			),
		)
	}
	if !addressSlicesEqual(e1.bcc, e2.bcc) {
		errors = append(
			errors,
			stackerrors.Newf(
				"`Bcc` does not match: expected %s, got %s",
				convertAddressListToReadable(e1.bcc),
				convertAddressListToReadable(e2.bcc),
			),
		)
	}
	if e1.subject != e2.subject {
		errors = append(errors, stackerrors.Newf("`subject` does not match: expected `%s`, got `%s`", e1.subject, e2.subject))
	}
	if e1.mainPart.alternativeSubPart.textPart.GetBody() != e2.mainPart.alternativeSubPart.textPart.GetBody() {
		errors = append(errors, stackerrors.Newf("`text` does not match: expected `%s`, got `%s`", e1.mainPart.alternativeSubPart.textPart.GetBody(), e2.mainPart.alternativeSubPart.textPart.GetBody()))
	}
	if e1.mainPart.alternativeSubPart.htmlPart.GetBody() != e2.mainPart.alternativeSubPart.htmlPart.GetBody() {
		errors = append(errors, stackerrors.Newf("`html` does not match: expected `%s`, got `%s`", e1.mainPart.alternativeSubPart.htmlPart.GetBody(), e2.mainPart.alternativeSubPart.htmlPart.GetBody()))
	}
	if len(e1.mainPart.embeddedSubParts.ExtractPartsSlice()) != len(e2.mainPart.embeddedSubParts.ExtractPartsSlice()) {
		errors = append(errors, stackerrors.Newf("`Embedded` count does not match: expected `%d`, got `%d`", len(e1.GetEmbedded().ExtractPartsSlice()), len(e2.GetEmbedded().ExtractPartsSlice())))
	} else {

		for ind, e1Embedded := range e1.mainPart.embeddedSubParts.ExtractPartsSlice() {
			if e1Embedded.GetBody() != e2.mainPart.embeddedSubParts.ExtractPartsSlice()[ind].GetBody() {
				errors = append(errors, stackerrors.Newf("`Embedded[%d]` body does not match: expected `%s`, got `%s`", ind, e1Embedded.GetBody(), e2.mainPart.embeddedSubParts.ExtractPartsSlice()[ind].GetBody()))
			}
		}
	}
	if len(e1.attachments.ExtractPartsSlice()) != len(e2.attachments.ExtractPartsSlice()) {
		errors = append(errors, stackerrors.Newf("`Attachments` count does not match: expected `%d`, got `%d`", len(e1.attachments.ExtractPartsSlice()), len(e2.attachments.ExtractPartsSlice())))
	} else {
		for ind, e1Attachment := range e1.attachments.ExtractPartsSlice() {
			if e1Attachment.GetBody() != e2.attachments.ExtractPartsSlice()[ind].GetBody() {
				errors = append(errors, stackerrors.Newf("`Attachments[%d]` body does not match: expected `%s`, got `%s`", ind, e1Attachment.GetBody(), e2.attachments.ExtractPartsSlice()[ind].GetBody()))
			}
		}
	}

	return errors
}

func convertAddressListToReadable(addresses []*mail.Address) string {
	addressValues := make([]string, 0)
	for _, addr := range addresses {
		addressValues = append(addressValues, fmt.Sprintf(`%#v`, addr))
	}
	return strings.Join(addressValues, `, `)
}
