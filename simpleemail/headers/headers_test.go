package headers

import (
	"reflect"
	"testing"
)

func TestEmptyHeader(t *testing.T) {
	h := NewHeaders().(*headers)
	if len(h.headers) > 0 {
		t.Error(`empty headers contains something`)
	}
	if h.GetFirstHeaderValue(`to`) != `` {
		t.Errorf(`"To" header in empty headers list contains value "%s"`, h.GetFirstHeaderValue(`to`))
	}
}

func TestWithHeader(t *testing.T) {
	testVal := `s`
	h := NewHeaders().WithHeader(`to`, testVal)
	if h.GetFirstHeaderValue(`to`) != testVal {
		t.Errorf(`headers.WithHeader() failure: header value expected to be "%s", got "%s"`, testVal, h.GetFirstHeaderValue(`to`))
	}
}

func TestWithoutHeader(t *testing.T) {
	testVal := `s`
	h := NewHeadersFromMap(map[string][]string{`to`: {testVal}})
	if h.WithoutHeader(`to`).GetFirstHeaderValue(`to`) != `` {
		t.Error(`headers.WithoutHeader() failure: header still not empty`)
	}
}

func TestWithAddedHeader(t *testing.T) {
	testVal := `a`
	h := NewHeaders().WithAddedHeader(`to`, testVal)
	h = h.WithAddedHeader(`to`, testVal+testVal)
	if len(h.GetHeaderValues(`to`)) != 2 {
		t.Error(`headers.WithAddedHeader() failure: header values length is not equal to 2`)
	}
	if h.GetHeaderValues(`To`)[0] != testVal {
		t.Errorf(`WithAddedHeader() failure: header added value expected to be "%s", got "%s"`, testVal, h.GetHeaderValues(`To`)[1])
	}
	if h.GetHeaderValues(`To`)[1] != testVal+testVal {
		t.Errorf(`headers.WithAddedHeader() failure: header added value expected to be "%s", got "%s"`, testVal+testVal, h.GetHeaderValues(`To`)[1])
	}
}

func TestCompile(t *testing.T) {
	h := NewHeaders().
		WithAddedHeader(`to`, `a`).
		WithAddedHeader(`to`, `b`).
		WithAddedHeader(`from`, `c`)

	expected := "From: c\r\nTo: a\r\nTo: b\r\n"

	got := string(h.Compile())
	if got != expected {
		t.Errorf(`headers.Compile() failure: expected to be "%s", got "%s"`, expected, got)
	}
}

func TestHeadersClone(t *testing.T) {
	h := NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	h2 := h.(*headers).clone()
	if h == h2 {
		t.Errorf(`headers.clone() failure: pointers are equal`)
	}
	if !reflect.DeepEqual(h.headers, h2.headers) {
		t.Errorf(`headers.clone() failure: inner maps differ`)
	}
}

func TestHeadersExtract(t *testing.T) {
	h := NewHeaders().WithAddedHeader(`to`, `s`).WithAddedHeader(`from`, `d`)
	h2 := h.ExtractHeadersMap()
	if !reflect.DeepEqual(map[string][]string(h.headers), h2) {
		t.Errorf(`headers.ExtractHeadersMap() failure: extracted map differ`)
	}
}

func TestHeadersFromMap(t *testing.T) {
	testMap := map[string][]string{`To`: {`a`}, `From`: {`b`}}
	h := NewHeadersFromMap(testMap)
	if !reflect.DeepEqual(map[string][]string(h.headers), testMap) {
		t.Errorf(`headers.NewHeadersFromMap() failure: imported map differ`)
	}
}
