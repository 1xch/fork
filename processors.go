package fork

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
)

type Processor interface {
	Widget
	Validater
	Filterer
	Errorer
}

type processor struct {
	Widget
	Errorer
	Validater
	Filterer
}

func NewProcessor(w Widget, validaters []interface{}, filters []interface{}) *processor {
	return &processor{
		Widget:    w,
		Errorer:   NewErrorer(),
		Validater: NewValidater(validaters...),
		Filterer:  NewFilterer(filters...),
	}
}

type Widget interface {
	String(interface{}) string
	Render(Field) template.HTML
	RenderWith(map[string]interface{}) template.HTML
}

const defaulttemplate = `
{{ define "fielderrors" }}<div class="field-errors"><ul>{{ range $x := .Errors }}<li>{{ $x }}</li>{{ end }}</ul></div>{{ end }}
{{ define "default" }}%s{{ if .Errors }}{{ template "fielderrors" .}}{{end}}{{ end }}
`

func NewWidget(t string) Widget {
	var err error
	ti := &widget{name: "default"}
	tt := fmt.Sprintf(defaulttemplate, t)
	ti.widget, err = template.New("widget").Parse(tt)
	if err != nil {
		ti.widget, _ = template.New("errorwidget").Parse(err.Error())
	}
	return ti
}

type widget struct {
	name   string
	widget *template.Template
}

func (w *widget) String(i interface{}) string {
	var buffer bytes.Buffer
	err := w.widget.ExecuteTemplate(&buffer, w.name, i)
	if err == nil {
		return buffer.String()
	}
	return err.Error()
}

func (w *widget) Render(f Field) template.HTML {
	return template.HTML(w.String(f))
}

func (w *widget) RenderWith(m map[string]interface{}) template.HTML {
	return template.HTML(w.String(m))
}

type Errorer interface {
	Errors(...string) []string
}

func NewErrorer() Errorer {
	return &errorer{}
}

type errorer struct {
	errors []string
}

func (e *errorer) Errors(ers ...string) []string {
	e.errors = append(e.errors, ers...)
	return e.errors
}

type Validater interface {
	Valid() bool
	Validate(Field) error
}

func NewValidater(v ...interface{}) Validater {
	return &validater{validaters: reflectValidaters(v...)}
}

type validater struct {
	valid      bool
	validated  bool
	validaters []reflect.Value
}

func (v *validater) Valid() bool {
	if !v.validated {
		return true
	}
	return v.valid
}

func (v *validater) Validate(f Field) error {
	for _, vdr := range v.validaters {
		err := Validate(vdr, f)
		if err != nil {
			v.validated, v.valid = true, false
		}
		return err
	}
	v.validated, v.valid = true, true
	return nil
}

func (v *validater) AddValidater(fn interface{}) {
	v.validaters = append(v.validaters, valueValidater(fn))
}

func reflectValidaters(fns ...interface{}) []reflect.Value {
	var ret []reflect.Value
	for _, fn := range fns {
		ret = append(ret, valueValidater(fn))
	}
	return ret
}

type Filterer interface {
	Filter(Field)
}

func NewFilterer(f ...interface{}) Filterer {
	return &filterer{filters: reflectFilters(f...)}
}

type filterer struct {
	filters []reflect.Value
}

func (fr *filterer) Filter(fd Field) {
	for _, fn := range fr.filters {
		_ = Filter(fn, fd)
	}
}

func (fr *filterer) AddFilter(fn interface{}) {
	fr.filters = append(fr.filters, valueFilter(fn))
}

func reflectFilters(fns ...interface{}) []reflect.Value {
	var ret []reflect.Value
	for _, fn := range fns {
		ret = append(ret, valueFilter(fn))
	}
	return ret
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func valueValidater(fn interface{}) reflect.Value {
	v := reflect.ValueOf(fn)
	if !isValidater(v.Type()) {
		panic(fmt.Sprintf("Cannot use %q as a validater function: function must return 1 value and the value must be an error", fn))
	}
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("%+v is not a function", fn))
	}
	return v
}

func valueFilter(fn interface{}) reflect.Value {
	v := reflect.ValueOf(fn)
	if !isFilter(v.Type()) {
		panic(fmt.Sprintf("Cannot use function %q with %d results\nreturn must be 1 value, or 1 value and 1 error value", fn, v.Type().NumOut()))
	}
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("%+v is not a function", fn))
	}
	return v
}

func isValidater(typ reflect.Type) bool {
	switch {
	case typ.NumOut() == 1 && typ.Out(0) == errorType:
		return true
	}
	return false
}

func isFilter(typ reflect.Type) bool {
	switch {
	case typ.NumOut() == 1:
		return true
	case typ.NumOut() == 2 && typ.Out(1) == errorType:
		return true
	}
	return false
}

func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}

func Validate(fn reflect.Value, args ...interface{}) error {
	validated, err := call(fn, args...)
	if err != nil {
		return err
	}
	if validated != nil {
		return validated.(error)
	}
	return nil
}

func Filter(fn reflect.Value, args ...interface{}) interface{} {
	filtered, err := call(fn, args...)
	if err != nil {
		return err
	}
	return filtered
}

func call(fn reflect.Value, args ...interface{}) (interface{}, error) {
	typ := fn.Type()
	numIn := typ.NumIn()
	var dddType reflect.Type
	if typ.IsVariadic() {
		if len(args) < numIn-1 {
			return nil, fmt.Errorf("wrong number of args: got %d want at least %d", len(args), numIn-1)
		}
		dddType = typ.In(numIn - 1).Elem()
	} else {
		if len(args) != numIn {
			return nil, fmt.Errorf("wrong number of args: got %d want %d", len(args), numIn)
		}
	}
	argv := make([]reflect.Value, len(args))
	for i, arg := range args {
		value := reflect.ValueOf(arg)
		// Compute the expected type. Clumsy because of variadics.
		var argType reflect.Type
		if !typ.IsVariadic() || i < numIn-1 {
			argType = typ.In(i)
		} else {
			argType = dddType
		}
		if !value.IsValid() && canBeNil(argType) {
			value = reflect.Zero(argType)
		}
		if !value.Type().AssignableTo(argType) {
			return nil, fmt.Errorf("arg %d has type %s; should be %s", i, value.Type(), argType)
		}
		argv[i] = value
	}
	result := fn.Call(argv)
	if len(result) == 2 && !result[1].IsNil() {
		return result[0].Interface(), result[1].Interface().(error)
	}
	return result[0].Interface(), nil
}
