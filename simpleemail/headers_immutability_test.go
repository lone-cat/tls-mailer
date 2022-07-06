package simpleemail

import (
	"reflect"
	"testing"
)

func TestWithHeaderImmutability(t *testing.T) {
	h := newHeaders().withHeader(`from`, `s`)
	h2 := h.clone()
	_ = h.withHeader(`to`, `s`)

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`withHeader() func immutability failure`)
	}
}

func TestWithoutHeaderImmutability(t *testing.T) {
	h := newHeaders().withHeader(`to`, `s`).withHeader(`from`, `s`)
	h2 := h.clone()
	_ = h.withoutHeader(`to`)

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`withoutHeader() func immutability failure`)
	}
}

func TestWithAddedHeaderImmutability(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `s`)
	h2 := h.clone()
	_ = h.withAddedHeader(`to`, `f`)

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`withAddedHeader() func immutability failure`)
	}
}

func TestNewHeadersFromMapImmutability(t *testing.T) {
	testMap := map[string][]string{`To`: {`a`}, `From`: {`b`}}
	h := newHeadersFromMap(testMap)
	h2 := h.clone()
	testMap[`To`][0] = `b`

	if !reflect.DeepEqual(h, h2) || reflect.DeepEqual(map[string][]string(h.headers), testMap) {
		t.Errorf(`newHeadersFromMap() func immutability failure`)
	}
}

func TestHeadersCloneImmutability(t *testing.T) {
	testMap := map[string][]string{`To`: {`a`}, `From`: {`b`}}
	h := newHeadersFromMap(testMap)
	h2 := h.clone()
	h2.headers[`To`][0] = `b`

	if reflect.DeepEqual(h, h2) {
		t.Errorf(`newHeadersFromMap() func immutability failure`)
	}
}

func TestExtractHeadersImmutability(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `s`)
	h2 := h.clone()
	headersMap := h.extractHeadersMap()
	headersMap[`To`][0] = `b`

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`extractHeadersMap() func immutability failure`)
	}
}

func TestGetHeaderValuesImmutability(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `s`)
	h2 := h.clone()
	headerSlice := h.getHeaderValues(`to`)
	headerSlice[0] = `b`

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`getHeaderValues() func immutability failure`)
	}
}
