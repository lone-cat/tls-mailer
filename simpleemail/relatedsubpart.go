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

func (p *relatedSubPart) WithText(text []byte) *relatedSubPart {
	newRelSubPart := p.clone()
	newRelSubPart.alternativeSubPart = newRelSubPart.alternativeSubPart.withText(text)
	return newRelSubPart
}

func (p *relatedSubPart) WithHtml(html []byte) *relatedSubPart {
	newRelSubPart := p.clone()
	newRelSubPart.alternativeSubPart = newRelSubPart.alternativeSubPart.withHtml(html)
	return newRelSubPart
}

func (p *relatedSubPart) WithAlternativeSubPart(altSubPart *alternativeSubPart) *relatedSubPart {
	newRelSubPart := p.clone()
	newRelSubPart.alternativeSubPart = altSubPart
	return newRelSubPart
}

func (p *relatedSubPart) WithEmbeddedSubPart(part part.Part) *relatedSubPart {
	newRelSubPart := p.clone()
	newRelSubPart.embeddedSubParts = newRelSubPart.embeddedSubParts.WithAppended(part)
	return newRelSubPart
}

func (p *relatedSubPart) WithoutEmbeddedSubParts() *relatedSubPart {
	newRelSubPart := p.clone()
	newRelSubPart.embeddedSubParts = part.NewPartsList()
	return newRelSubPart
}

func (p *relatedSubPart) Dump() map[string]any {
	if p == nil {
		return nil
	}
	dump := make(map[string]any)
	dump[`headers`] = p.headers.Dump()
	dump[`alternativeSubPart`] = p.alternativeSubPart.Dump()
	dump[`embeddedParts`] = p.embeddedSubParts.Dump()
	return dump
}

func (p *relatedSubPart) clone() *relatedSubPart {
	return &relatedSubPart{
		headers:            p.headers,
		alternativeSubPart: p.alternativeSubPart,
		embeddedSubParts:   p.embeddedSubParts,
	}
}

func (p *relatedSubPart) isEmpty() bool {
	return p.alternativeSubPart.isEmpty() && p.embeddedSubParts.IsEmpty()
}

func (p *relatedSubPart) toPart() part.Part {
	alternativePart := p.alternativeSubPart.toPart()
	if p.embeddedSubParts.IsEmpty() {
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
