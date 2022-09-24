package encode

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/quotedprintable"
	"strings"
	"unicode/utf8"
)

const mimeLineLength = 76

var base64LineSeparator = []byte("\r\n")

var Decoder = &mime.WordDecoder{}

func StringToBase64(val string) (result string, err error) {
	builder := &strings.Builder{}
	lineSplitter := NewSplitter(builder, base64LineSeparator, mimeLineLength)
	encoder := base64.NewEncoder(base64.StdEncoding, lineSplitter)

	valBytes := []byte(val)
	ln, err := encoder.Write(valBytes)
	if err != nil {
		return
	}
	if ln < len(valBytes) {
		return ``, errors.New(`written less bytes`)
	}
	err = encoder.Close()
	if err != nil {
		return
	}

	return builder.String(), nil
}

func BytesToBase64(b []byte) (res []byte, err error) {
	buffer := &bytes.Buffer{}
	lineSplitter := NewSplitter(buffer, base64LineSeparator, mimeLineLength)
	encoder := base64.NewEncoder(base64.StdEncoding, lineSplitter)

	ln, err := encoder.Write(b)
	if err != nil {
		return
	}
	if ln < len(b) {
		err = errors.New(`written less bytes`)
		return
	}
	err = encoder.Close()
	if err != nil {
		return
	}
	res = buffer.Bytes()
	return
}

func FromBase64(val string) (string, error) {
	bytess, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return ``, err
	}

	return string(bytess), nil
}

func StringToQuotedPrintable(s string) (string, error) {
	var ac bytes.Buffer
	w := quotedprintable.NewWriter(&ac)
	_, err := w.Write([]byte(s))
	if err != nil {
		return ``, err
	}
	err = w.Close()
	if err != nil {
		return ``, err
	}
	return ac.String(), nil
}

func BytesToQuotedPrintable(b []byte) (res []byte, err error) {
	var ac bytes.Buffer
	w := quotedprintable.NewWriter(&ac)
	_, err = w.Write(b)
	if err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	res = ac.Bytes()
	return
}

func FromQuotedPrintable(s string) (string, error) {
	r := quotedprintable.NewReader(strings.NewReader(s))
	result, err := io.ReadAll(r)
	if err != nil {
		return ``, err
	}
	res := strings.ReplaceAll(string(result), `_`, ` `)
	return res, nil
}

func EncodeHeader(headerValue string) string {
	return mime.QEncoding.Encode("utf-8", headerValue)
}

func EncodedHeaderToMultiline1(encodedHeader string) string {
	return strings.ReplaceAll(encodedHeader, `?= =?`, "?=\r\n =?")
}

func EncodedHeaderToMultiline(encodedHeader string) string {
	sourceParts := strings.Split(encodedHeader, ` `)
	resultLines := make([]string, 0)
	line := []string{sourceParts[0]}
	for _, subStr := range sourceParts[1:] {
		if utf8.RuneCountInString(strings.Join(line, ` `))+utf8.RuneCountInString(subStr) <= mimeLineLength {
			line = append(line, subStr)
		} else {
			resultLines = append(resultLines, strings.Join(line, ` `))
			line = []string{subStr}
		}
	}

	if len(line) > 0 {
		resultLines = append(resultLines, strings.Join(line, ` `))
	}

	return strings.Join(resultLines, "\r\n ")
}

func DecodeHeader(header string) (string, error) {
	return Decoder.DecodeHeader(header)
}
