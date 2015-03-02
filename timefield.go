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
		name:      name,
		format:    format,
		processor: NewProcessor(widget, validaters, filters),
	}
}

type timefield struct {
	name        string
	format      string
	Data        time.Time
	Information string
	*processor
}

func (t *timefield) New() Field {
	var newfield timefield = *t
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

func (t *timefield) Set(req *http.Request) {
	val := req.FormValue(t.Name())
	n, err := time.Parse(t.format, val)
	if err != nil {
		t.Errors(fmt.Sprintf("Cannot parse %s in format %s", val, t.format))
	}
	t.Data, t.Information = n, n.Format(t.format)
	t.Validate(t)
}

func datewidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<input type="date" name="{{ .Name }}" value="{{ .Information }}" %s>`, strings.Join(options, " ")))
}

func DateField(name string, options ...string) Field {
	return TimeField(name, dateFormat, datewidget(options...), nil, nil)
}
