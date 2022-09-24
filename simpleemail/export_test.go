package simpleemail

import "github.com/lone-cat/tls-mailer/simpleemail/part"

func (e *Email) GetAttachments() part.PartsList {
	return e.attachments
}

func (e *Email) GetEmbedded() part.PartsList {
	return e.mainPart.embeddedSubParts
}
