package simpleemail

func copySlice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

func copyHeadersMap(src map[string][]string) map[string][]string {
	dst := make(map[string][]string)
	for headerName, headerValuesSlice := range src {
		dst[headerName] = copySlice(headerValuesSlice)
	}
	return dst
}
