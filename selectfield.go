package fork

import (
	"fmt"
	"net/http"
	"strings"
)

type Selection struct {
	Value string
	Label string
	Set   bool
}

func NewSelection(value string, label string, set bool) *Selection {
	return &Selection{value, label, set}
}

type selectfield struct {
	Selections []*Selection
	*named
	Processor
}

func (s *selectfield) New(i ...interface{}) Field {
	var newfield selectfield = *s
	newfield.named = s.named.Copy()
	copy(newfield.Selections, s.Selections)
	s.SetValidateable(false)
	return &newfield
}

func (s *selectfield) Get() *Value {
	return NewValue(s.Selections)
}

func setselection(v string, selections []*Selection) {
	for _, s := range selections {
		if s.Value == v {
			s.Set = true
		}
	}
}

func (s *selectfield) Set(req *http.Request) {
	val := strings.Split(req.FormValue(s.Name()), " ")
	for _, v := range val {
		setselection(v, s.Selections)
	}
	s.SetValidateable(true)
}

const selectbase = `<select name="{{ .Name }}" %s>{{ range $x := .Selections }}<option value="{{ $x.Value }}"{{ if $x.Set }} selected{{ end }}>{{ $x.Label}}</option>{{ end }}</select>`

func selectWidget(options ...string) Widget {
	return NewWidget(WithOptions(selectbase, options...))
}

func SelectField(name string, s []*Selection, v []interface{}, f []interface{}, options ...string) Field {
	return &selectfield{
		Selections: s,
		named:      newnamed(name),
		Processor: NewProcessor(
			selectWidget(options...),
			NewValidater(v...),
			NewFilterer(f...),
		),
	}
}

type radiofield struct {
	Selections []Field
	*named
	Processor
}

func (r *radiofield) New(i ...interface{}) Field {
	var newfield radiofield = *r
	copy(newfield.Selections, r.Selections)
	r.SetValidateable(false)
	return &newfield
}

func (r *radiofield) Get() *Value {
	return NewValue(r.Selections)
}

func (r *radiofield) Set(req *http.Request) {
	for _, s := range r.Selections {
		s.Set(req)
	}
	r.SetValidateable(true)
}

func radioWidget(name string, legend string, options ...string) Widget {
	in := strings.Join([]string{
		fmt.Sprintf(`<fieldset name="%s" `, name),
		`%s>`,
		fmt.Sprintf(`<legend>%s</legend>`, legend),
		`<ul>{{ range $x := .Selections }}<li>{{ .Render $x }}</li>{{ end }}</ul></fieldset>`,
	}, "")
	return NewWidget(WithOptions(in, options...))
}

func makeradioinputs(name string, selections []*Selection) []Field {
	var ret []Field
	for _, s := range selections {
		ret = append(ret, RadioInput(name, s.Label, s.Value, false))
	}
	return ret
}

func RadioField(name string, legend string, s []*Selection, v []interface{}, f []interface{}, options ...string) Field {
	return &radiofield{
		Selections: makeradioinputs(name, s),
		named:      newnamed(name),
		Processor: NewProcessor(
			radioWidget(name, legend, options...),
			NewValidater(v...),
			NewFilterer(f...),
		),
	}
}
