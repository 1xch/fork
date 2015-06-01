package fork

import "net/http"

type Field interface {
	New() Field
	Fielder
}

type Fielder interface {
	Name(...string) string
	Get() *Value
	Set(*http.Request)
	Validateable() bool
	Processor
}

type Processor interface {
	Widget
	Validater
	Filterer
}

type processor struct {
	Widget
	Validater
	Filterer
}

func NewProcessor(w Widget, validaters []interface{}, filters []interface{}) *processor {
	return &processor{
		Widget:    w,
		Validater: NewValidater(validaters...),
		Filterer:  NewFilterer(filters...),
	}
}
