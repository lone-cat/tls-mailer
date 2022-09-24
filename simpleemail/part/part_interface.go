package part

import (
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"net/mail"
)

type Part interface {
	GetHeaders() headers.Headers
	WithHeaders(headers headers.Headers) Part
	GetBodyLen() int
	GetBody() []byte
	WithBody(body []byte) Part
	WithBodyFromString(body string) Part
	GetSubParts() []Part
	WithSubParts(subParts ...Part) Part
	Compile() ([]byte, error)
	ToPlainMessage() (*mail.Message, error)
	Clone() Part
}
