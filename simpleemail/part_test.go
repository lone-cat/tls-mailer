package simpleemail

import (
	"reflect"
	"testing"
)

func TestNewPart(t *testing.T) {
	p := newPart()
	if p.body != `` {
		t.Errorf(`newPart() failure: body is not empty`)
	}
	if len(p.headers.headers) != 0 {
		t.Errorf(`newPart() failure: headers are not empty`)
	}
}

func TestNewPartFromString(t *testing.T) {
	p := newPartFromString(text)
	if p.body != text {
		t.Errorf(`newPartFromString() failure: body expected "%s", got "%s"`, text, p.body)
	}
}

func TestPartClone(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `d`)
	p := newPart().
		withHeaders(h).
		withBody(text).
		withSubParts([]*part{{headers: h.clone(), body: text, subParts: newSubParts()}})

	p2 := p.clone()

	if p == p2 {
		t.Errorf(`part.clone() failure: pointers are equal`)
	}
	if !reflect.DeepEqual(p.headers, p2.headers) {
		t.Errorf(`part.clone() failure: header maps differ`)
	}
	if p.body != p2.body {
		t.Errorf(`part.clone() failure: bodies differ`)
	}
	if !reflect.DeepEqual(p.subParts, p2.subParts) {
		t.Errorf(`part.clone() failure: subParts differ`)
	}
}

func TestPartGetHeaders(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `d`)
	p := newPart().withHeaders(h)
	h2 := p.getHeaders()
	if !reflect.DeepEqual(p.headers, h2) {
		t.Errorf(`part.getHeaders() failure: headers map differ`)
	}
}

func TestPartWithHeaders(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `d`)
	p := newPart().withHeaders(h)
	if !reflect.DeepEqual(p.headers, h) {
		t.Errorf(`part.withHeaders() failure: headers map differ`)
	}
}

func TestPartGetBody(t *testing.T) {
	p := newPart()
	p.body = text
	if p.GetBody() != text {
		t.Errorf(`part.GetBody() failure: body differ`)
	}
}

func TestPartWithBody(t *testing.T) {
	p := newPart().withBody(text)
	if p.body != text {
		t.Errorf(`part.withBody() failure: body differ`)
	}
}

func TestPartGetSubparts(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `d`)
	subP := []*part{{headers: h.clone(), body: text, subParts: newSubParts()}}
	p := newPart()
	p.subParts = subP
	gotSubParts := p.getSubParts()
	if !reflect.DeepEqual(gotSubParts, subParts(subP)) {
		t.Errorf(`part.getSubParts() failure: subparts differ`)
	}
}

func TestPartWithSubparts(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `d`)
	subP := []*part{{headers: h.clone(), body: text, subParts: newSubParts()}}
	p := newPart().withSubParts(subP)
	if !reflect.DeepEqual(p.subParts, subParts(subP)) {
		t.Errorf(`part.withSubParts() failure: subparts differ`)
	}
}
