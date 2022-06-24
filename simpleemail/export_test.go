package simpleemail

func (e Email) GetAttachments() subParts {
	return e.attachments.clone()
}

func (e Email) GetEmbedded() subParts {
	return e.mainPart.embeddedSubParts.clone()
}
