package part

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/lone-cat/stackerrors"
	"github.com/lone-cat/tls-mailer/simpleemail/encoding"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

func ConvertMessageToPartRecursive(msg *mail.Message) (exportedPart Part, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get(headers.ContentType))

	if strings.HasPrefix(mediaType, headers.MultipartPrefix) {
		if err != nil {
			return
		}
		exportedPart, err = convertMultipartMsgToPart(msg, params[`boundary`])
	} else {
		exportedPart, err = convertSimpleMsgToPart(msg)
	}

	return
}

func convertSimpleMsgToPart(msg *mail.Message) (exportedPartAsInterface Part, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	exportedPart := &part{
		headers:  headers.NewHeadersFromMap(msg.Header),
		subParts: NewPartsList(),
	}

	var msgBodyBytes []byte
	msgBodyBytes, err = io.ReadAll(msg.Body)
	if err != nil {
		return
	}

	exportedPart.body = msgBodyBytes

	exportedPartAsInterface, err = unpackBody(exportedPart)

	return
}

func convertMultipartMsgToPart(msg *mail.Message, boundary string) (exportedPart Part, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	if boundary == `` {
		err = errors.New(`boundary is not set`)
		return
	}

	exportedPart = &part{
		headers:  headers.NewHeadersFromMap(msg.Header),
		subParts: NewPartsList(),
	}

	convertedSubParts := NewPartsList()
	mr := multipart.NewReader(msg.Body, boundary)
	var p *multipart.Part
	for {
		p, err = mr.NextRawPart()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		subMsg := &mail.Message{Header: mail.Header(p.Header), Body: p}
		var subPart Part
		subPart, err = ConvertMessageToPartRecursive(subMsg)
		if err != nil {
			return
		}
		convertedSubParts = convertedSubParts.WithAppended(subPart)
	}

	exportedPart = exportedPart.WithSubParts(convertedSubParts.ExtractPartsSlice()...)

	return
}

func unpackBody(prt Part) (unpacked Part, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()
	enc := prt.GetHeaders().GetContentTransferEncoding()

	unpacked = &part{
		headers:  prt.GetHeaders(),
		body:     prt.GetBody(),
		subParts: NewPartsList(prt.GetSubParts()...),
	}

	if enc == encoding.QuotedPrintable || enc == encoding.Base64 {
		var decodedBodyBytes []byte
		if enc == encoding.QuotedPrintable {
			decodedBodyBytes, err = io.ReadAll(quotedprintable.NewReader(bytes.NewReader(unpacked.GetBody())))
			if err != nil {
				return
			}
		}
		if enc == encoding.Base64 {
			decodedBodyBytes, err = io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(unpacked.GetBody())))
			if err != nil {
				return
			}
		}
		unpacked = unpacked.WithBody(decodedBodyBytes)
		unpacked = unpacked.WithHeaders(
			unpacked.GetHeaders().WithoutHeader(headers.ContentTransferEncoding),
		)
	}

	return
}
