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

const (
	ContentTransferEncodingHeader = `Content-Transfer-Encoding`
	ContentTypeHeader             = `Content-Type`
	FromHeader                    = `from`
	ToHeader                      = `to`
	CCHeader                      = `cc`
	BCCHeader                     = `bcc`
	SubjectHeader                 = `subject`
)

func (e Encoding) String() string {
	return string(e)
}

type Headers struct {
	headers mail.Header
}

func newHeaders() Headers {
	return Headers{
		headers: make(map[string][]string),
	}
}

func newHeadersFromMap(headers mail.Header) Headers {
	h := newHeaders()
	for headerName, headerValues := range headers {
		h = h.withHeader(headerName, headerValues...)
	}
	return h
}

func (h Headers) extractHeadersMap() map[string][]string {
	return copyHeadersMap(h.headers)
}

func (h Headers) clone() Headers {
	h.headers = copyHeadersMap(h.headers)
	return h
}

func (h Headers) withHeader(header string, values ...string) Headers {
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

func (h Headers) withAddedHeader(header string, values ...string) Headers {
	newHeaders := h.clone()
	if len(values) < 1 {
		return newHeaders
	}
	for _, val := range values {
		textproto.MIMEHeader(newHeaders.headers).Add(header, val)
	}
	return newHeaders
}

func (h Headers) withoutHeader(header string) Headers {
	newHeaders := h.clone()
	textproto.MIMEHeader(newHeaders.headers).Del(header)
	return newHeaders
}

func (h Headers) getFirstHeaderValue(header string) string {
	return h.headers.Get(header)
}

func (h Headers) getContentType() (contentType string, err error) {
	contentType, _, err = mime.ParseMediaType(h.getFirstHeaderValue(ContentTypeHeader))
	return
}

func (h Headers) getBoundary() (boundary string, err error) {
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

func (h Headers) isMultipartWithError() (bool, error) {
	contentType, err := h.getContentType()
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(contentType, MultipartPrefix), nil
}

func (h Headers) isMultipart() bool {
	multipart, _ := h.isMultipartWithError()
	return multipart
}

func (h Headers) getAddressList(header string) (addresses []mail.Address, err error) {
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

func (h Headers) getContentTransferEncoding() Encoding {
	return Encoding(strings.ToLower(h.getFirstHeaderValue(ContentTransferEncodingHeader)))
}

func (h Headers) compile() []byte {
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
