package simpleemail

type relatedSubPart struct {
	headers            Headers
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

func (m relatedSubPart) clone() relatedSubPart {
	newMainPart := newRelatedSubPart()
	newMainPart.headers = m.headers.clone()
	newMainPart.alternativeSubPart = m.alternativeSubPart.clone()
	newMainPart.embeddedSubParts = m.embeddedSubParts.clone()
	return newMainPart
}

func (m relatedSubPart) isEmpty() bool {
	return m.alternativeSubPart.isEmpty() && len(m.embeddedSubParts) < 1
}

func (m relatedSubPart) toPart() part {
	clonedParts := m.embeddedSubParts.clone()
	if !m.alternativeSubPart.isEmpty() {
		clonedParts = append([]part{m.alternativeSubPart.toPart()}, clonedParts...)
	}

	if len(clonedParts) < 1 {
		return newPart()
	}

	if len(clonedParts) == 1 {
		return clonedParts[0]
	}

	exportedPart := newPart()
	exportedPart.headers = m.headers.clone()
	exportedPart.subParts = clonedParts
	if !exportedPart.headers.isMultipart() {
		exportedPart.headers = exportedPart.headers.withHeader(`Content-Type`, MultipartRelated)
	}

	return exportedPart
}
