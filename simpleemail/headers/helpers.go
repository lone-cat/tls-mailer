package headers

import (
	"crypto/rand"
	"fmt"
	"github.com/lone-cat/tls-mailer/common"
	"io"
)

func copyHeadersMap(src map[string][]string) map[string][]string {
	dst := make(map[string][]string)
	for headerName, headerValuesSlice := range src {
		dst[headerName] = common.CloneSlice(headerValuesSlice)
	}
	return dst
}

// GenerateBoundary body for function is copied from mime/multipart/writer.go:randomBoundary() and slightly modified
func GenerateBoundary() string {
	var buf [25]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("boundary_%x", buf[:])
}
