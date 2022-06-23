package simpleemail

func (e Email) WithAttachedString(attachment string) Email {
	newEmail := e.clone()
	newEmail.attachments = append(newEmail.attachments, newAttachedPartFromString(attachment))
	return newEmail
}

func (e Email) WithEmbeddedString(embedded string) Email {
	newEmail := e.clone()
	newEmail.mainPart.embeddedSubParts = append(newEmail.mainPart.embeddedSubParts, newEmbeddedPartFromString(``, embedded))
	return newEmail
}

func (e Email) GetAttachments() subParts {
	return e.attachments.clone()
}

func (e Email) GetEmbedded() subParts {
	return e.mainPart.embeddedSubParts.clone()
}
