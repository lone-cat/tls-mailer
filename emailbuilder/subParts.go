package emailbuilder

type subParts []Part

func NewSubParts() subParts {
	return make([]Part, 0)
}

func (s subParts) Clone() subParts {
	newSubParts := make([]Part, 0, len(s))
	for _, subPart := range s {
		newSubParts = append(newSubParts, subPart.Clone())
	}
	return newSubParts
}

func generateBoundary() string {
	return `bound`
}
