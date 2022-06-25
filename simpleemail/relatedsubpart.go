package simpleemail

type relatedSubPart struct {
	headers            headers
	alternativeSubPart alternativeSubPart
	embeddedSubParts   subParts
}

func newRelatedSubPart() relatedSubPart {
	return relatedSubPart{
		headers:            newHeaders(),
		alternativeSubPart: newAlternativeSubPart(),
		embeddedSubParts:   newSubParts(),
	}
}

func (p relatedSubPart) clone() relatedSubPart {
	newMainPart := newRelatedSubPart()
	newMainPart.headers = p.headers.clone()
	newMainPart.alternativeSubPart = p.alternativeSubPart.clone()
	newMainPart.embeddedSubParts = p.embeddedSubParts.clone()
	return newMainPart
}

func (p relatedSubPart) isEmpty() bool {
	return p.alternativeSubPart.isEmpty() && len(p.embeddedSubParts) < 1
}

func (p relatedSubPart) toPart() part {
	alternativePart := p.alternativeSubPart.toPart()
	if len(p.embeddedSubParts) < 1 {
		return alternativePart
	}

	exportedPart := newPart()
	exportedPart.headers = p.headers.clone()
	exportedPart.subParts = append([]part{alternativePart}, p.embeddedSubParts...)
	if !exportedPart.headers.isMultipart() {
		exportedPart.headers = exportedPart.headers.withHeader(`Content-Type`, MultipartRelated)
	}

	return exportedPart
}
