package mailer

import (
	"errors"
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail"
	"net/mail"
	"net/smtp"
)

// golang enum! from here
type clientType uint8

const (
	undefinedClient clientType = iota
	TLSClient
	StartTLSClient
)

const (
	Undefined = `undefined`
	TLS       = `TLS`
	StartTLS  = `STARTTLS`
)

func (c clientType) String() string {
	switch c {
	case TLSClient:
		return TLS
	case StartTLSClient:
		return StartTLS
	default:
		return Undefined
	}
}

// to here

type Client struct {
	clientType clientType
	server     string
	auth       smtp.Auth
	sender     *mail.Address
}

func (c *Client) Send(email EmailForClient) error {
	err := getFirstBasicEmailValidationError(email, c.sender)
	if err != nil {
		return err
	}

	recipients, err := getValidRecipientsStringList(email.GetRecipients())
	if err != nil {
		return err
	}

	compiledBody, err := email.Compile()
	if err != nil {
		return err
	}

	err = validateEmailBody(compiledBody, c.sender, recipients)
	if err != nil {
		return err
	}

	return c.sendMail(c.server, c.auth, email.GetSender().Address, recipients, compiledBody)
}

func (c *Client) sendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	if c.clientType == TLSClient {
		return SendMail(addr, a, from, to, msg)
	} else if c.clientType == StartTLSClient {
		return smtp.SendMail(addr, a, from, to, msg)
	}
	return errors.New(`undefined client type`)
}

func (c *Client) GetType() string {
	return c.clientType.String()
}

func (c *Client) Email() simpleemail.Email {
	newMail := simpleemail.NewEmptyEmail()
	newMail, _ = newMail.WithFrom(c.sender)
	return newMail
}

func NewClient(
	clientType clientType,
	host string,
	port uint16,
	user string,
	password string,
	sender *mail.Address,
) (*Client, error) {
	if clientType.String() == Undefined {
		return nil, errors.New(`invalid client type`)
	}
	if sender == nil {
		return nil, errors.New(`nil sender passed to client constructor`)
	}
	if sender.Address == `` {
		return nil, errors.New(`empty sender passed to client constructor`)
	}

	auth := smtp.PlainAuth(``, user, password, host)
	cl := &Client{
		clientType: clientType,
		server:     fmt.Sprintf(`%s:%d`, host, port),
		auth:       auth,
		sender:     &mail.Address{Name: sender.Name, Address: sender.Address},
	}

	return cl, nil
}
