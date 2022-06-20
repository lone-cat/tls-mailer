package simpleemail

type mainSubPart struct {
	headers          Headers
	textSubPart      textSubPart
	embeddedSubParts subParts
}

func newMainSubPart() mainSubPart {
	return mainSubPart{
		headers:          newHeaders(),
		textSubPart:      newTextSubPart(),
		embeddedSubParts: newSubParts(),
	}
}

func (m mainSubPart) clone() mainSubPart {
	newMainPart := newMainSubPart()
	newMainPart.headers = m.headers.clone()
	newMainPart.textSubPart = m.textSubPart.clone()
	newMainPart.embeddedSubParts = m.embeddedSubParts.clone()
	return newMainPart
}

func (m mainSubPart) isEmpty() bool {
	return m.textSubPart.isEmpty() && len(m.embeddedSubParts) < 1
}

func (m mainSubPart) toPart() part {
	clonedParts := m.embeddedSubParts.clone()
	if !m.textSubPart.isEmpty() {
		clonedParts = append([]part{m.textSubPart.toPart()}, clonedParts...)
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
