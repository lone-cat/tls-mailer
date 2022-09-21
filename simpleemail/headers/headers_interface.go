package headers

import (
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
	GetContentTransferEncoding() Encoding
	Compile() []byte
	Date() (time.Time, error)
	AddressList(key string) ([]*mail.Address, error)
}
