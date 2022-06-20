package simpleemail

type textSubPart struct {
	headers  Headers
	textPart part
	htmlPart part
}

func newTextSubPart() textSubPart {
	return textSubPart{
		headers:  newHeaders(),
		textPart: newPart(),
		htmlPart: newPart(),
	}
}

func (t textSubPart) clone() textSubPart {
	clonedPart := newTextSubPart()
	clonedPart.headers = t.headers.clone()
	clonedPart.textPart = t.textPart.clone()
	clonedPart.htmlPart = t.htmlPart.clone()
	return clonedPart
}

func (t textSubPart) withText(text string) textSubPart {
	clonedPart := t.clone()
	clonedPart.textPart = clonedPart.textPart.withBody(text)
	return clonedPart
}

func (t textSubPart) withHtml(html string) textSubPart {
	clonedPart := t.clone()
	clonedPart.htmlPart = clonedPart.htmlPart.withBody(html)
	return clonedPart
}

func (t textSubPart) isTextEmpty() bool {
	return t.textPart.body == ``
}

func (t textSubPart) isHtmlEmpty() bool {
	return t.htmlPart.body == ``
}

func (t textSubPart) isEmpty() bool {
	return t.isTextEmpty() && t.isHtmlEmpty()
}

func (t textSubPart) toPart() part {
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

func (t textSubPart) compile() ([]byte, error) {
	return t.toPart().compile()
}
