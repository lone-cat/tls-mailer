package headers

import (
	"errors"
	"fmt"
	"github.com/lone-cat/tls-mailer/common"
	"github.com/lone-cat/tls-mailer/simpleemail/encoding"
	"mime"
	"net/mail"
	"net/textproto"
	"sort"
	"strings"
	"time"
)

const (
	ContentDisposition      = `Content-Disposition`
	ContentId               = `Content-Id`
	ContentTransferEncoding = `Content-Transfer-Encoding`
	ContentType             = `Content-Type`
	From                    = `From`
	To                      = `To`
	CC                      = `Cc`
	BCC                     = `Bcc`
	Subject                 = `Subject`
)

const (
	MultipartPrefix = `multipart/`

	MultipartMixed       = MultipartPrefix + `mixed`
	MultipartRelated     = MultipartPrefix + `related`
	MultipartAlternative = MultipartPrefix + `alternative`

	TextPlain = `text/plain`
	TextHtml  = `text/html`
)

type headers struct {
	headers mail.Header
}

func NewHeaders() Headers {
	return &headers{
		headers: make(map[string][]string),
	}
}

func NewHeadersFromMap(mailHeader mail.Header) Headers {
	head := make(map[string][]string)
	for headerName, headerValuesSlice := range mailHeader {
		head[textproto.CanonicalMIMEHeaderKey(headerName)] = common.CloneSlice(headerValuesSlice)
	}
	return &headers{
		headers: head,
	}
}

func (h *headers) ExtractHeadersMap() map[string][]string {
	return copyHeadersMap(h.headers)
}

func (h *headers) WithHeader(header string, values ...string) Headers {
	if len(values) < 1 {
		return h
	}

	newHeaders := h.ExtractHeadersMap()
	textproto.MIMEHeader(newHeaders).Set(header, values[0])
	for _, val := range values[1:] {
		textproto.MIMEHeader(newHeaders).Add(header, val)
	}

	return &headers{
		headers: newHeaders,
	}
}

func (h *headers) WithAddedHeader(header string, values ...string) Headers {
	if len(values) < 1 {
		return h
	}

	newHeaders := h.ExtractHeadersMap()
	for _, val := range values {
		textproto.MIMEHeader(newHeaders).Add(header, val)
	}

	return &headers{
		headers: newHeaders,
	}
}

func (h *headers) WithoutHeader(header string) Headers {
	newHeaders := h.ExtractHeadersMap()
	textproto.MIMEHeader(newHeaders).Del(header)

	return &headers{
		headers: newHeaders,
	}
}

func (h *headers) Clone() Headers {
	return NewHeadersFromMap(h.headers)
}

func (h *headers) GetFirstHeaderValue(header string) string {
	return h.headers.Get(header)
}

func (h *headers) GetHeaderValues(header string) []string {
	return common.CloneSlice(h.headers[textproto.CanonicalMIMEHeaderKey(header)])
}

func (h *headers) GetContentType() (contentType string, err error) {
	contentType, _, err = mime.ParseMediaType(h.GetFirstHeaderValue(ContentType))
	return
}

func (h *headers) GetBoundary() (boundary string, err error) {
	_, params, err := mime.ParseMediaType(h.GetFirstHeaderValue(ContentType))
	if err != nil {
		return ``, err
	}
	boundary, exists := params[`boundary`]
	if !exists || boundary == `` {
		return ``, errors.New(`boundary is not set`)
	}
	return
}

func (h *headers) IsMultipartWithError() (bool, error) {
	contentType, err := h.GetContentType()
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(contentType, MultipartPrefix), nil
}

func (h *headers) IsMultipart() bool {
	multipart, _ := h.IsMultipartWithError()
	return multipart
}

func (h *headers) GetContentTransferEncoding() encoding.Type {
	return encoding.Type(strings.ToLower(h.GetFirstHeaderValue(ContentTransferEncoding)))
}

func (h *headers) Date() (time.Time, error) {
	return h.headers.Date()
}

func (h *headers) AddressList(key string) ([]*mail.Address, error) {
	return h.headers.AddressList(key)
}

func (h *headers) Compile() []byte {
	headerNames := make([]string, 0)
	for k := range h.headers {
		headerNames = append(headerNames, k)
	}
	sort.Strings(headerNames)

	headerBytes := make([]byte, 0)
	for _, headerName := range headerNames {
		for _, headerValue := range h.headers[headerName] {
			headerLine := fmt.Sprintf("%s: %s\r\n", headerName, encoding.EncodedHeaderToMultiline(encoding.EncodeHeader(headerValue)))
			headerBytes = append(headerBytes, []byte(headerLine)...)
		}
	}

	return headerBytes
}
