package part

import (
	"github.com/lone-cat/tls-mailer/simpleemail"
	"github.com/lone-cat/tls-mailer/simpleemail/headers"
	"reflect"
	"testing"
)

func TestPartCloneImmutable(t *testing.T) {
	h := headers.NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	var p *part = NewPart().
		WithHeaders(h).
		WithBody(simpleemail.text).
		WithSubParts([]Part{nil}).(*part)

	p2 := p.Clone().(*part)

	p.body = `ss`
	p.headers.WithHeader(`to`, `ss`)
	p.WithSubParts(append(p2.GetSubParts().parts, &part{headers: h, body: simpleemail.text, subParts: newPartsList()}))

	if reflect.DeepEqual(p.headers, p2.headers) {
		t.Errorf(`part.clone() func immutability failure: headers equal`)
	}

	if p.body == p2.body {
		t.Errorf(`part.clone() func immutability failure: bodies equal`)
	}

	if reflect.DeepEqual(p.getSubParts(), p2.getSubParts()) {
		t.Errorf(`part.clone() func immutability failure: subparts equal`)
	}
}
