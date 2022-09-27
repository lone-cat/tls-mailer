package simpleemail

import (
	"fmt"
	"github.com/lone-cat/stackerrors"
	"net/mail"
	"reflect"
	"strings"
)

func emailsDiffErrors(e1 *Email, e2 *Email) []error {
	errors := make([]error, 0)
	if !addressSlicesEqual(e1.GetFrom(), e2.GetFrom()) {
		errors = append(
			errors,
			stackerrors.Newf(
				"`From` does not match: expected %s, got %s",
				convertAddressListToReadable(e1.GetFrom()),
				convertAddressListToReadable(e2.GetFrom()),
			),
		)
	}
	if !addressSlicesEqual(e1.GetTo(), e2.GetTo()) {
		errors = append(
			errors,
			stackerrors.Newf(
				"`To` does not match: expected %s, got %s",
				convertAddressListToReadable(e1.GetTo()),
				convertAddressListToReadable(e2.GetTo()),
			),
		)
	}
	if !addressSlicesEqual(e1.GetCc(), e2.GetCc()) {
		errors = append(
			errors,
			stackerrors.Newf(
				"`Cc` does not match: expected %s, got %s",
				convertAddressListToReadable(e1.GetCc()),
				convertAddressListToReadable(e2.GetCc()),
			),
		)
	}
	if !addressSlicesEqual(e1.GetBcc(), e2.GetBcc()) {
		errors = append(
			errors,
			stackerrors.Newf(
				"`Bcc` does not match: expected %s, got %s",
				convertAddressListToReadable(e1.GetBcc()),
				convertAddressListToReadable(e2.GetBcc()),
			),
		)
	}
	if !reflect.DeepEqual(e1.subject, e2.subject) {
		errors = append(errors, stackerrors.Newf("`subject` does not match: expected `%s`, got `%s`", e1.subject, e2.subject))
	}
	if !reflect.DeepEqual(e1.mainPart.alternativeSubPart.textPart.GetBody(), e2.mainPart.alternativeSubPart.textPart.GetBody()) {
		errors = append(errors, stackerrors.Newf("`text` does not match: expected `%s`, got `%s`", e1.mainPart.alternativeSubPart.textPart.GetBody(), e2.mainPart.alternativeSubPart.textPart.GetBody()))
	}
	if !reflect.DeepEqual(e1.mainPart.alternativeSubPart.htmlPart.GetBody(), e2.mainPart.alternativeSubPart.htmlPart.GetBody()) {
		errors = append(errors, stackerrors.Newf("`html` does not match: expected `%s`, got `%s`", e1.mainPart.alternativeSubPart.htmlPart.GetBody(), e2.mainPart.alternativeSubPart.htmlPart.GetBody()))
	}
	if len(e1.mainPart.embeddedSubParts.ExtractPartsSlice()) != len(e2.mainPart.embeddedSubParts.ExtractPartsSlice()) {
		errors = append(errors, stackerrors.Newf("`Embedded` count does not match: expected `%d`, got `%d`", len(e1.GetEmbedded().ExtractPartsSlice()), len(e2.GetEmbedded().ExtractPartsSlice())))
	} else {

		for ind, e1Embedded := range e1.mainPart.embeddedSubParts.ExtractPartsSlice() {
			if !reflect.DeepEqual(e1Embedded.GetBody(), e2.mainPart.embeddedSubParts.ExtractPartsSlice()[ind].GetBody()) {
				errors = append(errors, stackerrors.Newf("`Embedded[%d]` body does not match: expected `%s`, got `%s`", ind, e1Embedded.GetBody(), e2.mainPart.embeddedSubParts.ExtractPartsSlice()[ind].GetBody()))
			}
		}
	}
	if len(e1.attachments.ExtractPartsSlice()) != len(e2.attachments.ExtractPartsSlice()) {
		errors = append(errors, stackerrors.Newf("`Attachments` count does not match: expected `%d`, got `%d`", len(e1.attachments.ExtractPartsSlice()), len(e2.attachments.ExtractPartsSlice())))
	} else {
		for ind, e1Attachment := range e1.attachments.ExtractPartsSlice() {
			if !reflect.DeepEqual(e1Attachment.GetBody(), e2.attachments.ExtractPartsSlice()[ind].GetBody()) {
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
