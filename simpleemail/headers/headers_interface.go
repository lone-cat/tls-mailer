package headers

import (
	"github.com/lone-cat/tls-mailer/simpleemail/encoding"
	"net/mail"
	"time"
)

type Headers interface {
	ExtractHeadersMap() map[string][]string
	WithHeader(header string, values ...string) Headers
	WithAddedHeader(header string, values ...string) Headers
	WithoutHeader(header string) Headers
	Clone() Headers
	GetFirstHeaderValue(header string) string
	GetHeaderValues(header string) []string
	GetContentType() (string, error)
	GetBoundary() (boundary string, err error)
	IsMultipartWithError() (bool, error)
	IsMultipart() bool
	GetContentTransferEncoding() encoding.Type
	Compile() []byte
	Date() (time.Time, error)
	AddressList(key string) ([]*mail.Address, error)
}
