package simpleemail

import (
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"github.com/lone-cat/tls-mailer/simpleemail/part"
)

type relatedSubPart struct {
	headers            headers.Headers
	alternativeSubPart *alternativeSubPart
	embeddedSubParts   part.PartsList
}

func newRelatedSubPart() *relatedSubPart {
	return &relatedSubPart{
		headers:            headers.NewHeaders(),
		alternativeSubPart: newAlternativeSubPart(),
		embeddedSubParts:   part.NewPartsList(),
	}
}

func (p *relatedSubPart) clone() *relatedSubPart {
	return &relatedSubPart{
		headers:            p.headers.Clone(),
		alternativeSubPart: p.alternativeSubPart.clone(),
		embeddedSubParts:   part.NewPartsList(p.embeddedSubParts.ExtractPartsSlice()...),
	}
}

func (p *relatedSubPart) isEmpty() bool {
	return p.alternativeSubPart.isEmpty() && len(p.embeddedSubParts.ExtractPartsSlice()) < 1
}

func (p *relatedSubPart) toPart() part.Part {
	alternativePart := p.alternativeSubPart.toPart()
	if len(p.embeddedSubParts.ExtractPartsSlice()) < 1 {
		return alternativePart
	}

	exportedPart := part.NewPart().WithHeaders(p.headers).WithSubParts(append([]part.Part{alternativePart}, p.embeddedSubParts.ExtractPartsSlice()...)...)
	if !exportedPart.GetHeaders().IsMultipart() {
		exportedPart = exportedPart.WithHeaders(
			exportedPart.GetHeaders().WithHeader(`Content-Type`, headers.MultipartRelated),
		)
	}

	return exportedPart
}
