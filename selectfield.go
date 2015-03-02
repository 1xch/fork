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

var selectbase string = `<select name="{{ .Name }}" %s>{{ range $x := .Selections }}<option value="{{ $x.Value }}"{{ if $x.Set }} selected{{ end }}>{{ $x.Label}}</option>{{ end }}</select>`

func selectwidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(selectbase, strings.Join(options, " ")))
}

func SelectField(name string, s []*Selection, validaters []interface{}, filters []interface{}, options ...string) Field {
	return &selectfield{
		name:       name,
		Selections: s,
		processor:  NewProcessor(selectwidget(options...), validaters, filters),
	}
}

type selectfield struct {
	name       string
	Selections []*Selection
	*processor
}

func (s *selectfield) New() Field {
	var newfield selectfield = *s
	copy(newfield.Selections, s.Selections)
	return &newfield
}

func (s *selectfield) Name(name ...string) string {
	if len(name) > 0 {
		s.name = strings.Join(name, "-")
	}
	return s.name
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
}

type radiofield struct {
	name       string
	Selections []Field
	*processor
}

func (r *radiofield) New() Field {
	var newfield radiofield = *r
	copy(newfield.Selections, r.Selections)
	return &newfield
}

func (r *radiofield) Name(name ...string) string {
	if len(name) > 0 {
		r.name = strings.Join(name, "-")
	}
	return r.name
}

func (r *radiofield) Get() *Value {
	return NewValue(r.Selections)
}

func (r *radiofield) Set(req *http.Request) {
	for _, s := range r.Selections {
		s.Set(req)
	}
}

func radiowidget(name string, legend string, options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<fieldset name="%s" %s><legend>%s</legend><ul>{{ range $x := .Selections }}<li>{{ .Render $x }}</li>{{ end }}</ul></fieldset>`, name, strings.Join(options, " "), legend))
}

func makeradioinputs(name string, selections []*Selection) []Field {
	var ret []Field
	for _, s := range selections {
		ret = append(ret, RadioInput(name, s.Label, s.Value, false))
	}
	return ret
}

func RadioField(name string, legend string, s []*Selection, validaters []interface{}, filters []interface{}, options ...string) Field {
	return &radiofield{
		name:       name,
		Selections: makeradioinputs(name, s),
		processor:  NewProcessor(radiowidget(name, legend, options...), validaters, filters),
	}
}
