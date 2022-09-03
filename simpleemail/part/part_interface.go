package part

import (
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"net/mail"
)

type Part interface {
	GetHeaders() headers.Headers
	WithHeaders(headers headers.Headers) Part
	GetBody() string
	WithBody(body string) Part
	GetSubParts() *partsList
	WithSubParts(subParts *partsList) Part
	Compile() ([]byte, error)
	ToPlainMessage() (*mail.Message, error)
	Clone() Part
}
