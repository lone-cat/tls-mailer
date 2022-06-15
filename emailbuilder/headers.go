package emailbuilder

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

const (
	ContentTransferEncodingHeader = `Content-Transfer-Encoding`
	ContentTypeHeader             = `Content-Type`
	FromHeader                    = `From`
	ToHeader                      = `To`
	CCHeader                      = `Cc`
	BCCHeader                     = `Bcc`
	SubjectHeader                 = `Subject`
)

func (e Encoding) String() string {
	return string(e)
}

type Headers struct {
	headers mail.Header
}

func NewHeaders() Headers {
	return Headers{
		headers: make(map[string][]string),
	}
}

func NewHeadersFromMap(headers mail.Header) Headers {
	h := NewHeaders()
	for headerName, headerValues := range headers {
		h = h.WithHeader(headerName, headerValues...)
	}
	return h
}

func (h Headers) ExtractHeadersMap() map[string][]string {
	return copyHeadersMap(h.headers)
}

func (h Headers) Clone() Headers {
	h.headers = copyHeadersMap(h.headers)
	return h
}

func (h Headers) WithHeader(header string, values ...string) Headers {
	newHeaders := h.Clone()
	if len(values) < 1 {
		return newHeaders
	}
	textproto.MIMEHeader(newHeaders.headers).Set(header, values[0])
	for _, val := range values[1:] {
		textproto.MIMEHeader(newHeaders.headers).Add(header, val)
	}
	return newHeaders
}

func (h Headers) WithAddedHeader(header string, values ...string) Headers {
	newHeaders := h.Clone()
	if len(values) < 1 {
		return newHeaders
	}
	for _, val := range values {
		textproto.MIMEHeader(newHeaders.headers).Add(header, val)
	}
	return newHeaders
}

func (h Headers) WithoutHeader(header string) Headers {
	newHeaders := h.Clone()
	textproto.MIMEHeader(newHeaders.headers).Del(header)
	return newHeaders
}

func (h Headers) GetFirstHeaderValue(header string) string {
	return h.headers.Get(header)
}

func (h Headers) GetContentType() (contentType string, err error) {
	contentType, _, err = mime.ParseMediaType(h.GetFirstHeaderValue(ContentTypeHeader))
	return
}

func (h Headers) GetBoundary() (boundary string, err error) {
	_, params, err := mime.ParseMediaType(h.GetFirstHeaderValue(ContentTypeHeader))
	if err != nil {
		return ``, err
	}
	boundary, exists := params[`boundary`]
	if !exists || boundary == `` {
		return ``, errors.New(`boundary is not set`)
	}
	return
}

func (h Headers) IsMultipart() (bool, error) {
	contentType, err := h.GetContentType()
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(contentType, MultipartPrefix), nil
}

func (h Headers) GetAddressList(header string) (addresses []mail.Address, err error) {
	addresses = make([]mail.Address, 0)
	ptrs, err := h.headers.AddressList(header)
	if err != nil {
		return
	}
	for _, addressPointer := range ptrs {
		addresses = append(addresses, *addressPointer)
	}
	return
}

func (h Headers) GetContentTransferEncoding() Encoding {
	return Encoding(strings.ToLower(h.GetFirstHeaderValue(ContentTransferEncodingHeader)))
}

func (h Headers) Render() string {
	headerNames := make([]string, 0)
	for k := range h.headers {
		headerNames = append(headerNames, k)
	}
	sort.Strings(headerNames)

	headerLines := make([]string, 0)
	for _, headerName := range headerNames {
		for _, headerValue := range h.headers[headerName] {
			headerLine := fmt.Sprintf(`%s: %s`, headerName, EncodedHeaderToMultiline(EncodeHeader(headerValue)))
			headerLines = append(headerLines, headerLine)
		}
	}

	return strings.Join(headerLines, "\r\n") + "\r\n"
}

func EncodeHeader(headerValue string) string {
	return mime.QEncoding.Encode("utf-8", headerValue)
}

func EncodedHeaderToMultiline(encodedHeader string) string {
	return strings.ReplaceAll(encodedHeader, `?= =?`, "?=\r\n ?=")
}
