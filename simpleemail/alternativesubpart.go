package simpleemail

type alternativeSubPart struct {
	headers  Headers
	textPart part
	htmlPart part
}

func newAlternativeSubPart() alternativeSubPart {
	return alternativeSubPart{
		headers:  newHeaders(),
		textPart: newPart(),
		htmlPart: newPart(),
	}
}

func (t alternativeSubPart) clone() alternativeSubPart {
	clonedPart := newAlternativeSubPart()
	clonedPart.headers = t.headers.clone()
	clonedPart.textPart = t.textPart.clone()
	clonedPart.htmlPart = t.htmlPart.clone()
	return clonedPart
}

func (t alternativeSubPart) withText(text string) alternativeSubPart {
	clonedPart := t.clone()
	clonedPart.textPart = clonedPart.textPart.withBody(text)
	return clonedPart
}

func (t alternativeSubPart) withHtml(html string) alternativeSubPart {
	clonedPart := t.clone()
	clonedPart.htmlPart = clonedPart.htmlPart.withBody(html)
	return clonedPart
}

func (t alternativeSubPart) isTextEmpty() bool {
	return t.textPart.body == ``
}

func (t alternativeSubPart) isHtmlEmpty() bool {
	return t.htmlPart.body == ``
}

func (t alternativeSubPart) isEmpty() bool {
	return t.isTextEmpty() && t.isHtmlEmpty()
}

func (t alternativeSubPart) toPart() part {
	if t.isEmpty() {
		return newPart()
	}

	if t.isHtmlEmpty() {
		return t.textPart.clone()
	}

	if t.isTextEmpty() {
		return t.htmlPart.clone()
	}

	exportedPart := newPart()
	exportedPart.headers = t.headers.clone()
	exportedPart.subParts = []part{t.textPart.clone(), t.htmlPart.clone()}
	if !exportedPart.headers.isMultipart() {
		exportedPart.headers = exportedPart.headers.withHeader(`Content-Type`, MultipartAlternative)
	}

	return exportedPart
}

func (t alternativeSubPart) compile() ([]byte, error) {
	return t.toPart().compile()
}
