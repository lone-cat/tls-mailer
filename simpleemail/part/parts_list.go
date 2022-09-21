package part

type partsList struct {
	parts []Part
}

func NewPartsList(parts ...Part) PartsList {
	return &partsList{parts: copyPartsSlice(parts)}
}

func (l *partsList) ExtractPartsSlice() []Part {
	return copyPartsSlice(l.parts)
}

func (l *partsList) WithAppended(prt Part) PartsList {
	return NewPartsList(
		append(l.ExtractPartsSlice(), prt)...,
	)
}

func copyPartsSlice(parts []Part) []Part {
	partsSlice := make([]Part, len(parts))
	copy(partsSlice, parts)
	return partsSlice
}
