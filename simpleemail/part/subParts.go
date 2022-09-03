package part

type partsList struct {
	parts []Part
}

func newPartsList() *partsList {
	return &partsList{
		parts: make([]Part, 0),
	}
}

func (l *partsList) ExtractPartsSlice() []Part {
	return clonePartsSlice(l.parts)
}

func (l *partsList) WithAppended(prt Part) *partsList {
	return &partsList{
		parts: append(clonePartsSlice(l.parts), prt),
	}
}

func clonePartsSlice(src []Part) []Part {
	clonedSubParts := make([]Part, len(src))
	for index, subPart := range src {
		clonedSubParts[index] = subPart
	}

	return clonedSubParts
}
