package simpleemail

import (
	"github.com/lone-cat/stackerrors"
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

	var convertedPart *part
	convertedPart, err = convertMessageToPartRecursive(msg)
	if err != nil {
		return
	}

	email, err = convertPartToEmail(convertedPart)
	if err != nil {
		return
	}

	return
}
