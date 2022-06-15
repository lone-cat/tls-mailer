package emailbuilder

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/quotedprintable"
	"strings"
)

const base64LineLength = 76

var base64LineSeparator = []byte("\r\n")

var decoder = &mime.WordDecoder{}

func toBase64(val string) (result string, err error) {
	builder := &strings.Builder{}
	lineSplitter := NewSplitter(builder, base64LineSeparator, base64LineLength)
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

func fromBase64(val string) (string, error) {
	bytess, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return ``, err
	}

	return string(bytess), nil
}

func toQuotedPrintable(s string) (string, error) {
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
	//res := strings.ReplaceAll(ac.String(), ` `, `_`)
	return ac.String(), nil
}

func fromQuotedPrintable(s string) (string, error) {
	r := quotedprintable.NewReader(strings.NewReader(s))
	result, err := io.ReadAll(r)
	if err != nil {
		return ``, err
	}
	res := strings.ReplaceAll(string(result), `_`, ` `)
	return res, nil
}
