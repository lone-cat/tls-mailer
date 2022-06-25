package simpleemail

import (
	"net/textproto"
	"testing"
)

func TestWithHeaderImmutability(t *testing.T) {
	canonicalHeader := textproto.CanonicalMIMEHeaderKey(`to`)
	h := newHeaders()

	h2 := h.withHeader(canonicalHeader, `s`)

	if h2.getFirstHeaderValue(canonicalHeader) == h.getFirstHeaderValue(canonicalHeader) {
		t.Errorf(`asdasd`)
	}

	hMap := h.extractHeadersMap()
	hMap[canonicalHeader] = []string{`s`}

	if h2.getFirstHeaderValue(canonicalHeader) != hMap[`To`][0] {
		t.Errorf(`asdasd`)
	}
	if h.getFirstHeaderValue(canonicalHeader) != `` {
		t.Errorf(`asdasd`)
	}
}
