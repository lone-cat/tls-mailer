package simpleemail

import "fmt"

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

var i = 1

func generateBoundary() string {
	i++
	return fmt.Sprintf(`bound%d`, i)
}
