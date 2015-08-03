package fork

import (
	"net/http"
	"strings"
)

type Field interface {
	Named
	New(...interface{}) Field
	Get() *Value
	Set(*http.Request)
	Processor
}

type Named interface {
	Name() string
	ReName(...string) string
}

func newnamed(name string) *named {
	return &named{name: name}
}

type named struct {
	name string
}

func (n *named) Name() string {
	return n.name
}

func (n *named) ReName(rename ...string) string {
	if len(rename) > 0 {
		n.name = strings.Join(rename, "-")
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
	return &processor{w, v, f}
}
