package fork

import (
	"fmt"
	"net/http"
	"time"
)

const (
	dateFormat = "02/01/2006"
)

type timeField struct {
	format string
	Data   string
	*named
	Processor
}

func (t *timeField) New(i ...interface{}) Field {
	var newfield timeField = *t
	newfield.named = t.named.Copy()
	newfield.Data = ""
	newfield.SetValidateable(false)
	return &newfield
}

func (t *timeField) Get() *Value {
	return NewValue(t.Data)
}

func (t *timeField) Set(r *http.Request) {
	v := t.Filter(t.Name(), r)
	t.Data = v.String()
	t.SetValidateable(true)
}

func TimeField(name string, format string, widget Widget, v []interface{}, f []interface{}) Field {
	return &timeField{
		format: format,
		named:  newnamed(name),
		Processor: NewProcessor(
			widget,
			NewValidater(append(v, ValidateTime)...),
			NewFilterer(append(f, NewFilterTime(format))...),
		),
	}
}

func NewFilterTime(format string) func(string) string {
	return func(t string) string {
		n, err := time.Parse(format, t)
		if err == nil {
			return n.Format(format)
		}
		return t
	}
}

func ValidateTime(t *timeField) error {
	if t.Validateable() {
		_, err := time.Parse(t.format, t.Data)
		if err != nil {
			return fmt.Errorf("Cannot parse %s in format %s", t.Data, t.format)
		}
	}
	return nil
}

func dateWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="date" name="{{ .Name }}" value="{{ .Data }}" %s>`, options...))
}

func DateField(name string, options ...string) Field {
	return TimeField(name, dateFormat, dateWidget(options...), nil, nil)
}
