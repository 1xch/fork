package fork

import (
	"net/http"
	"strings"
)

type Field interface {
	New() Field
	Get() *Value
	Set(*http.Request)
	BaseField
	Processor
}

type BaseField interface {
	Name(...string) string
	Validateable() bool
}

func newBaseField(name string) *baseField {
	return &baseField{name: name}
}

type baseField struct {
	name         string
	validateable bool
}

func (b *baseField) Name(name ...string) string {
	if len(name) > 0 {
		b.name = strings.Join(name, "-")
	}
	return b.name
}

func (b *baseField) Validateable() bool {
	return b.validateable
}

func (b *baseField) Copy() *baseField {
	var ret baseField = *b
	return &ret
}

type Processor interface {
	Widget
	Filterer
	Validater
}

type processor struct {
	Widget
	Validater
	Filterer
}

func NewProcessor(w Widget, v Validater, f Filterer) *processor {
	return &processor{
		Widget:    w,
		Validater: v,
		Filterer:  f,
	}
}
