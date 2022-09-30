package simpleemail

import "github.com/lone-cat/tls-mailer/simpleemail/part"

func (e *email) GetAttachments() part.PartsList {
	return e.attachments
}

func (e *email) GetEmbedded() part.PartsList {
	return e.mainPart.embeddedSubParts
}
