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
	name         string
	Selection    *Selection
	validateable bool
	*processor
}

func (b *booleanfield) New() Field {
	var newfield booleanfield = *b
	var newselection Selection = *b.Selection
	newfield.Selection = &newselection
	b.validateable = false
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

func (b *booleanfield) Set(r *http.Request) {
	v := b.Filter(b.Name(), r)
	if v.String() == b.Selection.Value {
		b.Selection.Set = true
	} else {
		b.Selection.Set = false
	}
	b.validateable = true
}

func (b *booleanfield) Validateable() bool {
	return b.validateable
}

func togglewidget(input string, options ...string) Widget {
	in := strings.Join([]string{
		fmt.Sprintf(`<input type="%s" `, input),
		`name="{{ .Name }}" `,
		`value="{{ .Selection.Value }}"`,
		`{{ if .Selection.Set }} checked{{ end }} %s>`,
		`{{ .Selection.Label }}`,
	}, "")
	return NewWidget(WithOptions(in, options...))
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
