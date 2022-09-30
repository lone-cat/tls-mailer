package mailer

import "net/mail"

type EmailForClient interface {
	Compile() ([]byte, error)
	GetFrom() []*mail.Address
	GetTo() []*mail.Address
}
