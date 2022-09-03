package part

import (
	"github.com/lone-cat/tls-mailer/simpleemail"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"reflect"
	"testing"
)

func TestNewPart(t *testing.T) {
	p := NewPart().(*part)
	if p.body != `` {
		t.Errorf(`newPart() failure: body is not empty`)
	}
	if len(p.headers.ExtractHeadersMap()) != 0 {
		t.Errorf(`newPart() failure: headers are not empty`)
	}
}

func TestNewPartFromString(t *testing.T) {
	expected := `a`
	p := NewPartFromString(expected)
	actual := p.GetBody()
	if actual != expected {
		t.Errorf(`newPartFromString() failure: body expected "%s", got "%s"`, expected, actual)
	}
}

func TestPartClone(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	p := NewPart().
		WithHeaders(h).
		WithBody(simpleemail.text).
		WithSubParts([]*part{{headers: h.clone(), body: simpleemail.text, subParts: simpleemail.newSubParts()}})

	p2 := p.Clone()

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
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	p := newPart().withHeaders(h)
	h2 := p.getHeaders()
	if !reflect.DeepEqual(p.headers, h2) {
		t.Errorf(`part.getHeaders() failure: headers map differ`)
	}
}

func TestPartWithHeaders(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	p := newPart().withHeaders(h)
	if !reflect.DeepEqual(p.headers, h) {
		t.Errorf(`part.withHeaders() failure: headers map differ`)
	}
}

func TestPartGetBody(t *testing.T) {
	p := newPart()
	p.body = simpleemail.text
	if p.GetBody() != simpleemail.text {
		t.Errorf(`part.GetBody() failure: body differ`)
	}
}

func TestPartWithBody(t *testing.T) {
	p := newPart().withBody(simpleemail.text)
	if p.body != simpleemail.text {
		t.Errorf(`part.withBody() failure: body differ`)
	}
}

func TestPartGetSubparts(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	subP := []*part{{headers: h.clone(), body: simpleemail.text, subParts: simpleemail.newSubParts()}}
	p := newPart()
	p.subParts = subP
	gotSubParts := p.getSubParts()
	if !reflect.DeepEqual(gotSubParts, simpleemail.subParts(subP)) {
		t.Errorf(`part.getSubParts() failure: subparts differ`)
	}
}

func TestPartWithSubparts(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	subP := []*part{{headers: h.clone(), body: simpleemail.text, subParts: simpleemail.newSubParts()}}
	p := newPart().withSubParts(subP)
	if !reflect.DeepEqual(p.subParts, simpleemail.subParts(subP)) {
		t.Errorf(`part.withSubParts() failure: subparts differ`)
	}
}
