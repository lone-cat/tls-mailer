package part

import "github.com/lone-cat/tls-mailer/common"

type PartsList interface {
	ExtractPartsSlice() []Part
	WithAppended(prt Part) PartsList
	IsEmpty() bool
	Dump() []map[string]any
}

type partsList struct {
	parts []Part
}

func NewPartsList(parts ...Part) PartsList {
	return &partsList{parts: common.CloneSlice(parts)}
}

func (l *partsList) ExtractPartsSlice() []Part {
	return common.CloneSlice(l.parts)
}

func (l *partsList) WithAppended(prt Part) PartsList {
	return NewPartsList(
		append(l.ExtractPartsSlice(), prt)...,
	)
}

func (l *partsList) IsEmpty() bool {
	return len(l.parts) < 1
}

func (l *partsList) Dump() []map[string]any {
	if l == nil {
		return nil
	}
	dumpsSlice := make([]map[string]any, len(l.parts))
	for i := range l.parts {
		dumpsSlice[i] = l.parts[i].Dump()
	}

	return dumpsSlice
}
