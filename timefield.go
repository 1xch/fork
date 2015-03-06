package fork

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	dateFormat = "02/01/2006"
)

func TimeField(name string, format string, widget Widget, validaters []interface{}, filters []interface{}) Field {
	return &timefield{
		name:   name,
		format: format,
		processor: NewProcessor(
			widget,
			append(validaters, ValidateTime),
			append(filters, NewFilterTime(format)),
		),
	}
}

type timefield struct {
	name         string
	format       string
	Data         string
	validateable bool
	*processor
}

func (t *timefield) New() Field {
	var newfield timefield = *t
	t.validateable = false
	return &newfield
}

func (t *timefield) Name(name ...string) string {
	if len(name) > 0 {
		t.name = strings.Join(name, "-")
	}
	return t.name
}

func (t *timefield) Get() *Value {
	return NewValue(t.Data)
}

func (t *timefield) Set(r *http.Request) {
	v := t.Filter(t.Name(), r)
	t.Data = v.String()
	t.validateable = true
}

func (t *timefield) Validateable() bool {
	return t.validateable
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

func ValidateTime(t *timefield) error {
	if t.validateable {
		_, err := time.Parse(t.format, t.Data)
		if err != nil {
			return fmt.Errorf("Cannot parse %s in format %s", t.Data, t.format)
		}
	}
	return nil
}

func datewidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<input type="date" name="{{ .Name }}" value="{{ .Data }}" %s>`, strings.Join(options, " ")))
}

func DateField(name string, options ...string) Field {
	return TimeField(name, dateFormat, datewidget(options...), nil, nil)
}
