package fork

import (
	"fmt"
	"net/http"
	"strings"
)

type booleanfield struct {
	Selection *Selection
	*named
	Processor
}

func (b *booleanfield) New(fc ...FieldConfig) Field {
	var newfield booleanfield = *b
	var newselection Selection = *b.Selection
	newfield.named = b.named.Copy()
	newfield.Selection = &newselection
	newfield.SetValidateable(false)
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
	b.SetValidateable(true)
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
		named:     newnamed(name),
		Processor: NewProcessor(
			toggleWidget("checkbox", options...),
			nilValidater,
			nilFilterer,
		),
	}
}

func ToggleInput(name string, label string, value string, widget Widget, checked bool) Field {
	return &booleanfield{
		Selection: NewSelection(value, label, checked),
		named:     newnamed(name),
		Processor: NewProcessor(
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
