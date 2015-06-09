package fork

import (
	"net/http"
	"strings"
)

type Field interface {
	Named
	New() Field
	Get() *Value
	Set(*http.Request)
	Processor
}

type Named interface {
	Name(...string) string
}

func newnamed(name string) *named {
	return &named{name: name}
}

type named struct {
	name string
}

func (n *named) Name(name ...string) string {
	if len(name) > 0 {
		n.name = strings.Join(name, "-")
	}
	return n.name
}

func (n *named) Copy() *named {
	var ret named = *n
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
