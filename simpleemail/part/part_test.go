package part

import (
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"reflect"
	"testing"
)

func TestNewPart(t *testing.T) {
	p := NewPart().(*part)
	if len(p.body) > 0 {
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
	if string(actual) != expected {
		t.Errorf(`newPartFromString() failure: body expected "%s", got "%s"`, expected, actual)
	}
}

func TestPartClone(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	sp := NewPartsList(&part{headers: h.Clone(), body: []byte(text), subParts: NewPartsList()})
	p := NewPart().
		WithHeaders(h).
		WithBodyFromString(text).
		WithSubParts(sp.ExtractPartsSlice()...)

	p2 := p.Clone()

	if p == p2 {
		t.Errorf(`part.clone() failure: pointers are equal`)
	}
	if !reflect.DeepEqual(p.(*part).headers, p2.(*part).headers) {
		t.Errorf(`part.clone() failure: header maps differ`)
	}
	if !reflect.DeepEqual(p.(*part).body, p2.(*part).body) {
		t.Errorf(`part.clone() failure: bodies differ`)
	}
	if !reflect.DeepEqual(p.(*part).subParts, p2.(*part).subParts) {
		t.Errorf(`part.clone() failure: subParts differ`)
	}
}

func TestPartGetHeaders(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	p := NewPart().WithHeaders(h)
	h2 := p.GetHeaders()
	if !reflect.DeepEqual(p.(*part).headers, h2) {
		t.Errorf(`part.getHeaders() failure: headers map differ`)
	}
}

func TestPartWithHeaders(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	p := NewPart().WithHeaders(h)
	if !reflect.DeepEqual(p.(*part).headers, h) {
		t.Errorf(`part.withHeaders() failure: headers map differ`)
	}
}

func TestPartGetBody(t *testing.T) {
	p := NewPart()
	p.(*part).body = []byte(text)
	if string(p.GetBody()) != text {
		t.Errorf(`part.GetBody() failure: body differ`)
	}
}

func TestPartWithBody(t *testing.T) {
	p := NewPart().WithBodyFromString(text)
	if string(p.(*part).body) != text {
		t.Errorf(`part.withBody() failure: body differ`)
	}
}

func TestPartGetSubparts(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	subP := NewPartsList(&part{headers: h.Clone(), body: []byte(text), subParts: NewPartsList()})
	p := NewPart()
	p.(*part).subParts = subP
	gotSubParts := p.GetSubParts()
	if !reflect.DeepEqual(gotSubParts, subP.(*partsList).parts) {
		t.Errorf(`part.getSubParts() failure: subparts differ`)
	}
}

func TestPartWithSubparts(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	subP := NewPartsList(&part{headers: h.Clone(), body: []byte(text), subParts: NewPartsList()})
	p := NewPart().WithSubParts(subP.ExtractPartsSlice()...)
	if !reflect.DeepEqual(p.(*part).subParts, subP) {
		t.Errorf(`part.withSubParts() failure: subparts differ`)
	}
}
