package simpleemail

import (
	"reflect"
	"testing"
)

func TestPartCloneImmutable(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `d`)
	p := newPart().
		withHeaders(h).
		withBody(text).
		withSubParts([]*part{{headers: h.clone(), body: text, subParts: newSubParts()}})

	p2 := p.clone()

	p.body = `ss`
	p.headers.headers[`To`][0] = `ss`
	p.subParts[0] = &part{headers: h.clone(), body: text + `s`, subParts: newSubParts()}

	if reflect.DeepEqual(p.headers, p2.headers) {
		t.Errorf(`part.clone() func immutability failure: headers equal`)
	}

	if p.body == p2.body {
		t.Errorf(`part.clone() func immutability failure: bodies equal`)
	}

	if reflect.DeepEqual(p.subParts, p2.subParts) {
		t.Errorf(`part.clone() func immutability failure: subparts equal`)
	}
}
