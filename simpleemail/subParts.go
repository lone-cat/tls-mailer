package simpleemail

type subParts []part

func newSubParts() subParts {
	return make([]part, 0)
}

func (s subParts) clone() subParts {
	clonedSubParts := make([]part, len(s))
	for index, subPart := range s {
		clonedSubParts[index] = subPart.clone()
	}
	return clonedSubParts
}

func generateBoundary() string {
	return `bound`
}
