package simpleemail

import (
	"encoding/base64"
	"errors"
	"github.com/lone-cat/stackerrors"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

func convertMessageToPartRecursive(msg *mail.Message) (exportedPart *part, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get(ContentTypeHeader))

	if strings.HasPrefix(mediaType, MultipartPrefix) {
		if err != nil {
			return
		}
		exportedPart, err = convertMultipartMsgToPart(msg, params[`boundary`])
	} else {
		exportedPart, err = convertSimpleMsgToPart(msg)
	}

	return
}

func convertSimpleMsgToPart(msg *mail.Message) (exportedPart *part, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	exportedPart = &part{
		headers:  newHeadersFromMap(msg.Header),
		subParts: newSubParts(),
	}

	var msgBodyBytes []byte
	msgBodyBytes, err = io.ReadAll(msg.Body)
	if err != nil {
		return
	}

	exportedPart.body = string(msgBodyBytes)

	exportedPart, err = unpackBody(exportedPart)

	return
}

func convertMultipartMsgToPart(msg *mail.Message, boundary string) (exportedPart *part, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()

	if boundary == `` {
		err = errors.New(`boundary is not set`)
		return
	}

	exportedPart = &part{
		headers:  newHeadersFromMap(msg.Header),
		subParts: newSubParts(),
	}

	convertedSubParts := newSubParts()
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
		var subPart *part
		subPart, err = convertMessageToPartRecursive(subMsg)
		if err != nil {
			return
		}
		convertedSubParts = append(convertedSubParts, subPart)
	}

	exportedPart = exportedPart.withSubParts(convertedSubParts)

	return
}

func unpackBody(part *part) (unpacked *part, err error) {
	defer func() {
		err = stackerrors.WrapInDefer(err)
	}()
	encoding := part.getHeaders().getContentTransferEncoding()

	unpacked = part.clone()
	if encoding == EncodingQuotedPrintable || encoding == EncodingBase64 {
		var decodedBodyBytes []byte
		if encoding == EncodingQuotedPrintable {
			decodedBodyBytes, err = io.ReadAll(quotedprintable.NewReader(strings.NewReader(unpacked.body)))
			if err != nil {
				return
			}
		}
		if encoding == EncodingBase64 {
			decodedBodyBytes, err = io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(unpacked.body)))
			if err != nil {
				return
			}
		}
		unpacked.body = string(decodedBodyBytes)
		unpacked.headers = unpacked.headers.withoutHeader(ContentTransferEncodingHeader)
	}

	return
}
