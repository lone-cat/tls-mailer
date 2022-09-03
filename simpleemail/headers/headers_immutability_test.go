package headers

import (
	"reflect"
	"testing"
)

func TestWithHeaderImmutability(t *testing.T) {
	h := NewHeaders().WithHeader(`from`, `s`)
	h2 := h.clone()
	_ = h.WithHeader(`to`, `s`)

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`WithHeader() func immutability failure`)
	}
}

func TestWithoutHeaderImmutability(t *testing.T) {
	h := NewHeaders().WithHeader(`to`, `s`).WithHeader(`from`, `s`)
	h2 := h.clone()
	_ = h.WithoutHeader(`to`)

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`WithoutHeader() func immutability failure`)
	}
}

func TestWithAddedHeaderImmutability(t *testing.T) {
	h := NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `s`)
	h2 := h.clone()
	_ = h.WithAddedHeader(`to`, `f`)

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`WithAddedHeader() func immutability failure`)
	}
}

func TestNewHeadersFromMapImmutability(t *testing.T) {
	testMap := map[string][]string{`To`: {`a`}, `From`: {`b`}}
	h := NewHeadersFromMap(testMap)
	h2 := h.clone()
	testMap[`To`][0] = `b`

	if !reflect.DeepEqual(h, h2) || reflect.DeepEqual(map[string][]string(h.headers), testMap) {
		t.Errorf(`NewHeadersFromMap() func immutability failure`)
	}
}

func TestHeadersCloneImmutability(t *testing.T) {
	testMap := map[string][]string{`To`: {`a`}, `From`: {`b`}}
	h := NewHeadersFromMap(testMap)
	h2 := h.clone()
	h2.headers[`To`][0] = `b`

	if reflect.DeepEqual(h, h2) {
		t.Errorf(`NewHeadersFromMap() func immutability failure`)
	}
}

func TestExtractHeadersImmutability(t *testing.T) {
	h := NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `s`)
	h2 := &(*(h.(*headers)))
	headersMap := h.ExtractHeadersMap()
	headersMap[`To`][0] = `b`

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`ExtractHeadersMap() func immutability failure`)
	}
}

func TestGetHeaderValuesImmutability(t *testing.T) {
	h := NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `s`)
	h2 := h.clone()
	headerSlice := h.GetHeaderValues(`to`)
	headerSlice[0] = `b`

	if !reflect.DeepEqual(h, h2) {
		t.Errorf(`GetHeaderValues() func immutability failure`)
	}
}
