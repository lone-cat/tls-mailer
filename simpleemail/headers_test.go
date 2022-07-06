package simpleemail

import (
	"reflect"
	"testing"
)

func TestEmptyHeader(t *testing.T) {
	h := newHeaders()
	if len(h.headers) > 0 {
		t.Error(`empty headers contains something`)
	}
	if h.getFirstHeaderValue(`to`) != `` {
		t.Errorf(`"To" header in empty headers list contains value "%s"`, h.getFirstHeaderValue(`to`))
	}
}

func TestWithHeader(t *testing.T) {
	testVal := `s`
	h := newHeaders().withHeader(`to`, testVal)
	if h.getFirstHeaderValue(`to`) != testVal {
		t.Errorf(`headers.withHeader() failure: header value expected to be "%s", got "%s"`, testVal, h.getFirstHeaderValue(`to`))
	}
}

func TestWithoutHeader(t *testing.T) {
	testVal := `s`
	h := newHeadersFromMap(map[string][]string{`to`: {testVal}})
	if h.withoutHeader(`to`).getFirstHeaderValue(`to`) != `` {
		t.Error(`headers.withoutHeader() failure: header still not empty`)
	}
}

func TestWithAddedHeader(t *testing.T) {
	testVal := `a`
	h := newHeaders().withAddedHeader(`to`, testVal)
	h = h.withAddedHeader(`to`, testVal+testVal)
	if len(h.getHeaderValues(`to`)) != 2 {
		t.Error(`headers.withAddedHeader() failure: header values length is not equal to 2`)
	}
	if h.getHeaderValues(`To`)[0] != testVal {
		t.Errorf(`withAddedHeader() failure: header added value expected to be "%s", got "%s"`, testVal, h.getHeaderValues(`To`)[1])
	}
	if h.getHeaderValues(`To`)[1] != testVal+testVal {
		t.Errorf(`headers.withAddedHeader() failure: header added value expected to be "%s", got "%s"`, testVal+testVal, h.getHeaderValues(`To`)[1])
	}
}

func TestCompile(t *testing.T) {
	h := newHeaders().
		withAddedHeader(`to`, `a`).
		withAddedHeader(`to`, `b`).
		withAddedHeader(`from`, `c`)

	expected := "From: c\r\nTo: a\r\nTo: b\r\n"

	got := string(h.compile())
	if got != expected {
		t.Errorf(`headers.compile() failure: expected to be "%s", got "%s"`, expected, got)
	}
}

func TestHeadersClone(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `d`)
	h2 := h.clone()
	if h == h2 {
		t.Errorf(`headers.clone() failure: pointers are equal`)
	}
	if !reflect.DeepEqual(h.headers, h2.headers) {
		t.Errorf(`headers.clone() failure: inner maps differ`)
	}
}

func TestHeadersExtract(t *testing.T) {
	h := newHeaders().withAddedHeader(`to`, `s`).withAddedHeader(`from`, `d`)
	h2 := h.extractHeadersMap()
	if !reflect.DeepEqual(map[string][]string(h.headers), h2) {
		t.Errorf(`headers.extractHeadersMap() failure: extracted map differ`)
	}
}

func TestHeadersFromMap(t *testing.T) {
	testMap := map[string][]string{`To`: {`a`}, `From`: {`b`}}
	h := newHeadersFromMap(testMap)
	if !reflect.DeepEqual(map[string][]string(h.headers), testMap) {
		t.Errorf(`headers.newHeadersFromMap() failure: imported map differ`)
	}
}
