package fork

import (
	"fmt"
	"net/http"
	"strings"
)

func BooleanField(name string, label string, start bool, options ...string) Field {
	return &booleanfield{
		name:      name,
		Selection: NewSelection(name, label, start),
		processor: NewProcessor(togglewidget("checkbox", options...), nil, nil),
	}
}

type booleanfield struct {
	name      string
	Selection *Selection
	*processor
}

func (b *booleanfield) New() Field {
	var newfield booleanfield = *b
	var newselection Selection = *b.Selection
	newfield.Selection = &newselection
	return &newfield
}

func (b *booleanfield) Name(name ...string) string {
	if len(name) > 0 {
		b.name = strings.Join(name, "-")
	}
	return b.name
}

func (b *booleanfield) Get() *Value {
	return NewValue(b.Selection)
}

func (b *booleanfield) Set(req *http.Request) {
	val := req.FormValue(b.Name())
	if val == b.Selection.Value {
		b.Selection.Set = true
	} else {
		b.Selection.Set = false
	}
}

func togglewidget(input string, options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<input type="%s" name="{{ .Name }}" value="{{ .Selection.Value }}"{{ if .Selection.Set }} checked{{ end }} %s>{{ .Selection.Label }}`, input, strings.Join(options, " ")))
}

func ToggleInput(name string, label string, value string, widget Widget, checked bool) Field {
	return &booleanfield{
		name:      name,
		Selection: NewSelection(value, label, checked),
		processor: NewProcessor(widget, nil, nil),
	}
}

func RadioInput(name string, label string, value string, checked bool, options ...string) Field {
	return ToggleInput(name, label, value, togglewidget("radio", options...), checked)
}

func CheckboxInput(name string, label string, value string, checked bool, options ...string) Field {
	return ToggleInput(name, label, value, togglewidget("checkbox", options...), checked)
}
