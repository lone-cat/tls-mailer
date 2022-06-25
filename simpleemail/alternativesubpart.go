package simpleemail

type alternativeSubPart struct {
	headers  headers
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

func (p alternativeSubPart) clone() alternativeSubPart {
	clonedPart := newAlternativeSubPart()
	clonedPart.headers = p.headers.clone()
	clonedPart.textPart = p.textPart.clone()
	clonedPart.htmlPart = p.htmlPart.clone()
	return clonedPart
}

func (p alternativeSubPart) withText(text string) alternativeSubPart {
	clonedPart := p.clone()
	clonedPart.textPart = clonedPart.textPart.withBody(text)
	return clonedPart
}

func (p alternativeSubPart) withHtml(html string) alternativeSubPart {
	clonedPart := p.clone()
	clonedPart.htmlPart = clonedPart.htmlPart.withBody(html)
	return clonedPart
}

func (p alternativeSubPart) isTextEmpty() bool {
	return p.textPart.body == ``
}

func (p alternativeSubPart) isHtmlEmpty() bool {
	return p.htmlPart.body == ``
}

func (p alternativeSubPart) isEmpty() bool {
	return p.isTextEmpty() && p.isHtmlEmpty()
}

func (p alternativeSubPart) toPart() part {
	if p.isEmpty() {
		return newPart()
	}

	if p.isHtmlEmpty() {
		return p.textPart.clone()
	}

	if p.isTextEmpty() {
		return p.htmlPart.clone()
	}

	exportedPart := newPart()
	exportedPart.headers = p.headers.clone()
	exportedPart.subParts = []part{p.textPart.clone(), p.htmlPart.clone()}
	if !exportedPart.headers.isMultipart() {
		exportedPart.headers = exportedPart.headers.withHeader(`Content-Type`, MultipartAlternative)
	}

	return exportedPart
}

func (p alternativeSubPart) compile() ([]byte, error) {
	return p.toPart().compile()
}
