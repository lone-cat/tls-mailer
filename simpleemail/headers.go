package simpleemail

import (
	"errors"
	"fmt"
	"mime"
	"net/mail"
	"net/textproto"
	"sort"
	"strings"
)

type Encoding string

const (
	EncodingEmpty           Encoding = ``
	Encoding7bit            Encoding = `7bit`
	Encoding8bit            Encoding = `8bit`
	EncodingBinary          Encoding = `binary`
	EncodingQuotedPrintable Encoding = `quoted-printable`
	EncodingBase64          Encoding = `base64`
)

func (e Encoding) String() string {
	return string(e)
}

const (
	ContentDispositionHeader      = `Content-Disposition`
	ContentIdHeader               = `Content-Id`
	ContentTransferEncodingHeader = `Content-Transfer-Encoding`
	ContentTypeHeader             = `Content-Type`
	FromHeader                    = `From`
	ToHeader                      = `To`
	CCHeader                      = `Cc`
	BCCHeader                     = `Bcc`
	SubjectHeader                 = `Subject`
)

type headers struct {
	headers mail.Header
}

func newHeaders() *headers {
	return &headers{
		headers: make(map[string][]string),
	}
}

func newHeadersFromMap(mailHeader mail.Header) *headers {
	return &headers{
		headers: copyHeadersMap(mailHeader),
	}
}

func (h *headers) extractHeadersMap() map[string][]string {
	return copyHeadersMap(h.headers)
}

func (h *headers) clone() *headers {
	return &headers{
		headers: copyHeadersMap(h.headers),
	}
}

func (h *headers) withHeader(header string, values ...string) *headers {
	newHeaders := h.clone()
	if len(values) < 1 {
		return newHeaders
	}
	textproto.MIMEHeader(newHeaders.headers).Set(header, values[0])
	for _, val := range values[1:] {
		textproto.MIMEHeader(newHeaders.headers).Add(header, val)
	}
	return newHeaders
}

func (h *headers) withAddedHeader(header string, values ...string) *headers {
	newHeaders := h.clone()
	if len(values) < 1 {
		return newHeaders
	}
	for _, val := range values {
		textproto.MIMEHeader(newHeaders.headers).Add(header, val)
	}
	return newHeaders
}

func (h *headers) withoutHeader(header string) *headers {
	newHeaders := h.clone()
	textproto.MIMEHeader(newHeaders.headers).Del(header)
	return newHeaders
}

func (h *headers) getFirstHeaderValue(header string) string {
	return h.headers.Get(header)
}

func (h *headers) getContentType() (contentType string, err error) {
	contentType, _, err = mime.ParseMediaType(h.getFirstHeaderValue(ContentTypeHeader))
	return
}

func (h *headers) getBoundary() (boundary string, err error) {
	_, params, err := mime.ParseMediaType(h.getFirstHeaderValue(ContentTypeHeader))
	if err != nil {
		return ``, err
	}
	boundary, exists := params[`boundary`]
	if !exists || boundary == `` {
		return ``, errors.New(`boundary is not set`)
	}
	return
}

func (h *headers) isMultipartWithError() (bool, error) {
	contentType, err := h.getContentType()
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(contentType, MultipartPrefix), nil
}

func (h *headers) isMultipart() bool {
	multipart, _ := h.isMultipartWithError()
	return multipart
}

func (h *headers) getContentTransferEncoding() Encoding {
	return Encoding(strings.ToLower(h.getFirstHeaderValue(ContentTransferEncodingHeader)))
}

func (h *headers) compile() []byte {
	headerNames := make([]string, 0)
	for k := range h.headers {
		headerNames = append(headerNames, k)
	}
	sort.Strings(headerNames)

	headerBytes := make([]byte, 0)
	for _, headerName := range headerNames {
		for _, headerValue := range h.headers[headerName] {
			headerLine := fmt.Sprintf("%s: %s\r\n", headerName, encodedHeaderToMultiline(encodeHeader(headerValue)))
			headerBytes = append(headerBytes, []byte(headerLine)...)
		}
	}

	return headerBytes
}
