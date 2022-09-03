package simpleemail

import (
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"github.com/lone-cat/tls-mailer/simpleemail/part"
)

type relatedSubPart struct {
	headers            headers.Headers
	alternativeSubPart *alternativeSubPart
	embeddedSubParts   part.subParts
}

func newRelatedSubPart() *relatedSubPart {
	return &relatedSubPart{
		headers:            headers.NewHeaders(),
		alternativeSubPart: newAlternativeSubPart(),
		embeddedSubParts:   part.newSubParts(),
	}
}

func (p *relatedSubPart) clone() *relatedSubPart {
	return &relatedSubPart{
		headers:            p.headers.clone(),
		alternativeSubPart: p.alternativeSubPart.clone(),
		embeddedSubParts:   p.embeddedSubParts.clone(),
	}
}

func (p *relatedSubPart) isEmpty() bool {
	return p.alternativeSubPart.isEmpty() && len(p.embeddedSubParts) < 1
}

func (p *relatedSubPart) toPart() *part.part {
	alternativePart := p.alternativeSubPart.toPart()
	if len(p.embeddedSubParts) < 1 {
		return alternativePart
	}

	exportedPart := &part.part{
		headers:  p.headers.clone(),
		subParts: append([]*part.part{alternativePart}, p.embeddedSubParts...),
	}
	if !exportedPart.headers.IsMultipart() {
		exportedPart.headers = exportedPart.headers.WithHeader(`Content-Type`, part.MultipartRelated)
	}

	return exportedPart
}
