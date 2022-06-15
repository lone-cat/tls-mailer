package tls_mailer

/*
import (
	"errors"
	"fmt"
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

func (c clientType) String() string {
	switch c {
	case TLSClient:
		return `TLS`
	case StartTLSClient:
		return `STARTTLS`
	default:
		return `undefined`
	}
}

// to here

type Client struct {
	clientType clientType
	server     string
	auth       smtp.Auth
	sender     *mail.Address
}

func (c *Client) Send(email *Email) error {
	compiledBody, err := email.compile()
	if err != nil {
		return err
	}
	to := make([]string, 0)
	for _, addr := range email.to {
		to = append(to, addr.Address)
	}
	return c.sendMail(c.server, c.auth, email.from.Address, to, []byte(compiledBody))
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

func (c *Client) Email(subject string, body string, to ...*mail.Address) (*Email, error) {
	realTo := make([]*mail.Address, 0)

	for _, addr := range to {
		if addr != nil {
			realTo = append(realTo, addr)
		}
	}

	if len(realTo) < 1 {
		return nil, errors.New(`no one reciever passed`)
	}

	return &Email{
		headers:  newHeaders(),
		from:     c.sender,
		to:       realTo,
		Subject:  subject,
		subParts: make([]*part, 0),
	}, nil
}

func NewClient(
	clientType clientType,
	host string,
	port uint16,
	user string,
	password string,
	sender *mail.Address,
) (*Client, error) {
	auth := smtp.PlainAuth("", user, password, host)
	cl := &Client{
		clientType: clientType,
		server:     fmt.Sprintf(`%s:%d`, host, port),
		auth:       auth,
		sender:     sender,
	}

	return cl, nil
}
*/
