package part

import "github.com/lone-cat/tls-mailer/common"

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
