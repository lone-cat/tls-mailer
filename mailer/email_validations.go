package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"net/mail"
)

// pre-send validations for email message

func getFirstBasicEmailValidationError(email EmailForClient, clientSender *mail.Address) error {
	if email == nil {
		return errors.New(`nil email passed`)
	}
	senders := email.GetFrom()
	if len(senders) < 1 {
		return errors.New(`email sender is not set`)
	}

	sender := senders[0]
	if sender == nil {
		return errors.New(`email sender is not set`)
	}

	if clientSender == nil {
		return errors.New(`client sender is not set`)
	}

	if sender.Address != clientSender.Address {
		return errors.New(`email sender does not match client sender`)
	}

	return nil
}

func getValidRecipientsStringList(rawRecipients []*mail.Address) ([]string, error) {
	recipients := make([]string, 0)
	for _, addr := range rawRecipients {
		if addr != nil {
			addrStr := addr.Address
			if addrStr != `` {
				recipients = append(recipients, addrStr)
			}
		}
	}
	if len(recipients) < 1 {
		return nil, errors.New(`email recipients does not contain valid values`)
	}

	return recipients, nil
}

func validateEmailBody(msgBytes []byte, clientSender *mail.Address, emailRecipientStrings []string) error {
	if len(msgBytes) < 1 {
		return errors.New(`email body is empty`)
	}
	if clientSender == nil {
		return errors.New(`client sender is nil`)
	}
	if len(emailRecipientStrings) < 1 {
		return errors.New(`empty recipients list passed`)
	}

	reader := bytes.NewReader(msgBytes)
	msg, err := mail.ReadMessage(reader)
	if err != nil {
		return err
	}

	err = validateFromHeader(msg.Header, clientSender)
	if err != nil {
		return err
	}

	err = validateToHeader(msg.Header, emailRecipientStrings)
	if err != nil {
		return err
	}

	return nil
}

func validateFromHeader(headers mail.Header, clientSender *mail.Address) error {
	senders, err := headers.AddressList(`from`)
	if err != nil {
		return err
	}

	if len(senders) != 1 {
		return errors.New(fmt.Sprintf(`message can be sent only from exactly one sender. senders passed: %d`, len(senders)))
	}

	sender := senders[0]
	if sender == nil {
		return errors.New(`message sender is nil`)
	}

	if sender.Address != clientSender.Address {
		return errors.New(`message sender does not match client sender`)
	}

	return nil
}

func validateToHeader(headers mail.Header, emailRecipientStrings []string) error {
	recipients, err := headers.AddressList(`to`)
	if err != nil {
		return err
	}
	if len(recipients) < 1 {
		return errors.New(`recipients must be included at least in "to" header`)
	}
	emailRecipientsMap := make(map[string]struct{})
	for _, v := range emailRecipientStrings {
		emailRecipientsMap[v] = struct{}{}
	}

	for _, recipient := range recipients {
		if recipient == nil {
			return errors.New(`message recipient is nil`)
		}
		_, exists := emailRecipientsMap[recipient.Address]
		if !exists {
			return errors.New(fmt.Sprintf(`message recipient "%s" from "to" header is not included in recipients list`, recipient.Address))
		}
	}

	return nil
}
