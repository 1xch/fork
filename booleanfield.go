package fork

import (
	"fmt"
	"net/http"
	"strings"
)

type booleanfield struct {
	Selection *Selection
	*baseField
	*processor
}

func (b *booleanfield) New() Field {
	var newfield booleanfield = *b
	var newselection Selection = *b.Selection
	newfield.baseField = b.baseField.Copy()
	newfield.Selection = &newselection
	b.validateable = false
	return &newfield
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

func toggleWidget(input string, options ...string) Widget {
	in := strings.Join([]string{
		fmt.Sprintf(`<input type="%s" `, input),
		`name="{{ .Name }}" `,
		`value="{{ .Selection.Value }}"`,
		`{{ if .Selection.Set }} checked{{ end }} %s>`,
		`{{ .Selection.Label }}`,
	}, "")
	return NewWidget(WithOptions(in, options...))
}

func BooleanField(name string, label string, start bool, options ...string) Field {
	return &booleanfield{
		Selection: NewSelection(name, label, start),
		baseField: newBaseField(name),
		processor: NewProcessor(
			toggleWidget("checkbox", options...),
			nilValidater,
			nilFilterer,
		),
	}
}

func ToggleInput(name string, label string, value string, widget Widget, checked bool) Field {
	return &booleanfield{
		Selection: NewSelection(value, label, checked),
		baseField: newBaseField(name),
		processor: NewProcessor(
			widget,
			nilValidater,
			nilFilterer,
		),
	}
}

func RadioInput(name string, label string, value string, checked bool, options ...string) Field {
	return ToggleInput(name, label, value, toggleWidget("radio", options...), checked)
}

func CheckboxInput(name string, label string, value string, checked bool, options ...string) Field {
	return ToggleInput(name, label, value, toggleWidget("checkbox", options...), checked)
}
