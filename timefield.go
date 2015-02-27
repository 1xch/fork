package fork

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	dateFormat         = "02/01/2006"
	defaultdatePattern = `(?:(?:0[1-9]|1[0-2])[\/\\-. ]?(?:0[1-9]|[12][0-9])|(?:(?:0[13-9]|1[0-2])[\/\\-. ]?30)|(?:(?:0[13578]|1[02])[\/\\-. ]?31))[\/\\-. ]?(?:19|20)[0-9]{2}`
)

func datewidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<input type="date" pattern="%s" title="Date as DD/MM/YYYY" name="{{ .Name }}" value="{{ .Information }}" %s>`, defaultdatePattern, strings.Join(options, " ")))
}

func TimeField(name string, format string, widget Widget) Field {
	return &timefield{
		name:      name,
		format:    format,
		processor: NewProcessor(widget, nil, nil),
	}
}

func DateField(name string, options ...string) Field {
	return TimeField(name, dateFormat, datewidget(options...))
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
	t.Information = val
	newtime, err := time.Parse(t.format, t.Information)
	if err != nil {
		t.Information = err.Error()
	}
	t.Data = newtime
}
