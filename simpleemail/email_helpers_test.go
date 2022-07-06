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
	if e1.mainPart.alternativeSubPart.textPart.body != e2.mainPart.alternativeSubPart.textPart.body {
		errors = append(errors, stackerrors.Newf("`text` does not match: expected `%s`, got `%s`", e1.mainPart.alternativeSubPart.textPart.body, e2.mainPart.alternativeSubPart.textPart.body))
	}
	if e1.mainPart.alternativeSubPart.htmlPart.body != e2.mainPart.alternativeSubPart.htmlPart.body {
		errors = append(errors, stackerrors.Newf("`html` does not match: expected `%s`, got `%s`", e1.mainPart.alternativeSubPart.htmlPart.body, e2.mainPart.alternativeSubPart.htmlPart.body))
	}
	if len(e1.mainPart.embeddedSubParts) != len(e2.mainPart.embeddedSubParts) {
		errors = append(errors, stackerrors.Newf("`Embedded` count does not match: expected `%d`, got `%d`", len(e1.GetEmbedded()), len(e2.GetEmbedded())))
	} else {
		for ind, e1Embedded := range e1.mainPart.embeddedSubParts {
			if e1Embedded.body != e2.mainPart.embeddedSubParts[ind].body {
				errors = append(errors, stackerrors.Newf("`Embedded[%d]` body does not match: expected `%s`, got `%s`", ind, e1Embedded.body, e2.mainPart.embeddedSubParts[ind].body))
			}
		}
	}
	if len(e1.attachments) != len(e2.attachments) {
		errors = append(errors, stackerrors.Newf("`Attachments` count does not match: expected `%d`, got `%d`", len(e1.attachments), len(e2.attachments)))
	} else {
		for ind, e1Attachment := range e1.attachments {
			if e1Attachment.body != e2.attachments[ind].body {
				errors = append(errors, stackerrors.Newf("`Attachments[%d]` body does not match: expected `%s`, got `%s`", ind, e1Attachment.body, e2.attachments[ind].body))
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
