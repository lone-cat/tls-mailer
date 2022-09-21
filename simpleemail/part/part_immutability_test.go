package part

import (
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"reflect"
	"testing"
)

func TestPartCloneImmutable(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	p := NewPart().
		WithHeaders(h).
		WithBody(text).
		WithSubParts([]Part{nil}...).(*part)

	p2 := p.Clone().(*part)

	p.body = `ss`

	p.headers = p.headers.WithHeader(`to`, `ss`)

	p.subParts = NewPartsList(append(p2.subParts.ExtractPartsSlice(), &part{headers: h, body: text, subParts: NewPartsList()})...)

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
