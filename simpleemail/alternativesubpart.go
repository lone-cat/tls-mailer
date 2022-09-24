package simpleemail

import (
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"github.com/lone-cat/tls-mailer/simpleemail/part"
)

type alternativeSubPart struct {
	headers  headers.Headers
	textPart part.Part
	htmlPart part.Part
}

func newAlternativeSubPart() *alternativeSubPart {
	return &alternativeSubPart{
		headers:  headers.NewHeaders(),
		textPart: part.NewPart(),
		htmlPart: part.NewPart(),
	}
}

func (p *alternativeSubPart) clone() *alternativeSubPart {
	return &alternativeSubPart{
		headers:  p.headers,
		textPart: p.textPart,
		htmlPart: p.htmlPart,
	}
}

func (p *alternativeSubPart) withText(text string) *alternativeSubPart {
	clonedPart := p.clone()
	clonedPart.textPart = clonedPart.textPart.WithBody(text)
	return clonedPart
}

func (p *alternativeSubPart) withHtml(html string) *alternativeSubPart {
	clonedPart := p.clone()
	clonedPart.htmlPart = clonedPart.htmlPart.WithBody(html)
	return clonedPart
}

func (p *alternativeSubPart) isTextEmpty() bool {
	return p.textPart.GetBody() == ``
}

func (p *alternativeSubPart) isHtmlEmpty() bool {
	return p.htmlPart.GetBody() == ``
}

func (p *alternativeSubPart) isEmpty() bool {
	return p.isTextEmpty() && p.isHtmlEmpty()
}

func (p *alternativeSubPart) toPart() part.Part {
	if p.isEmpty() {
		return part.NewPart()
	}

	if p.isHtmlEmpty() {
		return p.textPart
	}

	if p.isTextEmpty() {
		return p.htmlPart
	}

	exportedPart := part.NewPart()
	exportedPart = exportedPart.WithHeaders(p.headers)
	if !exportedPart.GetHeaders().IsMultipart() {
		exportedPart = exportedPart.WithHeaders(
			exportedPart.GetHeaders().
				WithHeader(
					`Content-Type`,
					headers.MultipartAlternative,
				),
		)
	}

	exportedPart = exportedPart.WithSubParts(p.textPart, p.htmlPart)

	return exportedPart
}

func (p *alternativeSubPart) compile() ([]byte, error) {
	return p.toPart().Compile()
}
