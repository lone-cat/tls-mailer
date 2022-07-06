package simpleemail

type relatedSubPart struct {
	headers            *headers
	alternativeSubPart *alternativeSubPart
	embeddedSubParts   subParts
}

func newRelatedSubPart() *relatedSubPart {
	return &relatedSubPart{
		headers:            newHeaders(),
		alternativeSubPart: newAlternativeSubPart(),
		embeddedSubParts:   newSubParts(),
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

func (p *relatedSubPart) toPart() *part {
	alternativePart := p.alternativeSubPart.toPart()
	if len(p.embeddedSubParts) < 1 {
		return alternativePart
	}

	exportedPart := &part{
		headers:  p.headers.clone(),
		subParts: append([]*part{alternativePart}, p.embeddedSubParts...),
	}
	if !exportedPart.headers.isMultipart() {
		exportedPart.headers = exportedPart.headers.withHeader(`Content-Type`, MultipartRelated)
	}

	return exportedPart
}
