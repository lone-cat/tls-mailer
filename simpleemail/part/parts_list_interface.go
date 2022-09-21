package part

type PartsList interface {
	ExtractPartsSlice() []Part
	WithAppended(prt Part) PartsList
}
