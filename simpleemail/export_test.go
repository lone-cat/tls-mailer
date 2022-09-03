package simpleemail

import "github.com/lone-cat/tls-mailer/simpleemail/part"

func (e *Email) GetAttachments() part.subParts {
	return e.attachments.clone()
}

func (e *Email) GetEmbedded() part.subParts {
	return e.mainPart.embeddedSubParts.clone()
}
