package simpleemail

import (
	"github.com/lone-cat/stackerrors"
	"github.com/lone-cat/tls-mailer/simpleemail/part"
	"net/mail"
	"strings"
)

func Import(message string) (email *Email, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	var msg *mail.Message
	msg, err = mail.ReadMessage(strings.NewReader(message))
	if err != nil {
		return
	}

	var convertedPart *part.part
	convertedPart, err = part.convertMessageToPartRecursive(msg)
	if err != nil {
		return
	}

	email, err = convertPartToEmail(convertedPart)
	if err != nil {
		return
	}

	return
}
