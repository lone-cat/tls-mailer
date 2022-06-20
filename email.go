package tls_mailer

import "net/mail"

type EmailForClient interface {
	Compile() ([]byte, error)
	GetSender() *mail.Address
	GetRecipients() []*mail.Address
}
