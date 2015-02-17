package fork

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"
	"text/template"
)

func newform(fs ...Fielder) Former {
	return &form{fields: fs}
}

type form struct {
	fields []Fielder
}

func (f *form) New() Former {
	var newform form = *f
	//fs := f.fields
	//newform.fields = nil
	//for _, field := range fs {
	//	newform.fields = append(newform.fields, field.New())
	//}
	return &newform
}

func (f *form) Fields(fs ...Fielder) []Fielder {
	f.fields = append(f.fields, fs...)
	return f.fields
}

func (f *form) Process(r *http.Request) {
	for _, fd := range f.Fields() {
		fd.Set(r.FormValue(fd.Name()))
	}
}

func (f *form) Valid() bool {
	for _, fd := range f.Fields() {
		if !fd.Valid() {
			return false
		}
	}
	return true
}

func (f *form) Render() string {
	b := new(bytes.Buffer)
	for _, fd := range f.Fields() {
		_, _ = b.WriteString(fd.Render())
	}
	return b.String()
}

func newfield(name string, widget Widget) *field {
	return &field{
		name:   name,
		widget: widget,
	}
}

type field struct {
	name       string
	data       *Value
	widget     Widget
	valid      bool
	validated  bool
	validaters []reflect.Value
	filters    []reflect.Value
	errors     []string
}

func (f *field) New() Fielder {
	var newfield field = *f
	return &newfield
}

func (f *field) Name() string {
	return f.name
}

func (f *field) Get() *Value {
	return f.data
}

func (f *field) Set(i interface{}) {
	f.data = NewValue(i)
}

func (f *field) Render() string {
	return f.widget.Render(f)
}

func (f *field) Valid() bool {
	if !f.validated {
		f.validate()
		f.validated = true
	}
	if len(f.errors) < 1 {
		return true
	}
	return false
}

func (f *field) validate() {
	for _, v := range f.validaters {
		err := validate(v, f.data.Raw) //v(f)
		if err != nil {
			f.errors = append(f.errors, err.Error())
		}
	}
}

func (f *field) Validaters(fns ...interface{}) []reflect.Value {
	for _, validater := range vs {
		v := reflect.ValueOf(validater)
		if v.Kind() != reflect.Func {
			panic(errors.New("Provided:(%+v, type: %T), but it is not a function", fn, fn))
		}
		if !isValidater(v.Type()) {
			panic(errors.New("Cannot use function %q with %d results\nreturn must be 1 value, or 1 value and 1 error value", fn, v.Type().NumOut()))
		}
		f.validaters = append(f.validaters, v)
	}
	return f.validaters
}

//func (f *field) filter() {
//	for _, filter := range f.filters {
//		filter(f, f.data)
//	}
//}

//func (f *field) Filters(fs ...Filter) []Filter {
//	for _, fr := range fs {
//		f.filters = append(f.filters, fr)
//	}
//	return f.filters
//}

func (f *field) Errors(e ...string) []string {
	f.errors = append(f.errors, e...)
	return f.errors
}

type TestWidget struct {
	widget *template.Template
}

func (t *TestWidget) Render(f Field) string {
	var buffer bytes.Buffer
	err := t.widget.Execute(&buffer, f)
	if err == nil {
		return buffer.String()
	}
	return err.Error()
}

func NewTestWidget(t string) *TestWidget {
	var err error
	ti := &TestWidget{}
	ti.widget, err = template.New("widget").Parse(t)
	if err != nil {
		ti.widget, _ = template.New("errorwidget").Parse(err.Error())
	}
	return ti
}

var (
	f1 Fielder = newfield("fruit", NewTestWidget(`<input type="text" name="{{ .Name }}" value="{{ .Get }}"></input>`))
	f2 Fielder = newfield("vegetable", NewTestWidget(`<input type="text" name="{{ .Name }}" value="{{ .Get }}"></input>`))
	f3 Fielder = newfield("submit", NewTestWidget(`<input type="submit" value="Submit"></input>`))
)

//var AcceptableFruits = []string{"apples, pears, cherries"}

//func isacceptable(fruit string) bool {
//	for _, f := range AcceptableFruits {
//		if fruit == f {
//			return true
//		}
//	}
//	return false
//}

//func AcceptableFruit(f Fielder) error {
//	for _, fr := range f.Data().([]interface{}) {
//		if !isacceptable(fr.(string)) {
//			return errors.New(fmt.Sprintf("%s is not an acceptable fruit", fr))
//		}
//	}
//	return nil
//}

//func UnacceptableVegetable(f Fielder) error {
//for _, vg := range f.Data().([]string) {
//	if vg == "spinach" {
//		return nil
//	}
//}
//return errors.New("This vegetable is unacceptable.")
//}

//func Isstring(f Fielder, i interface{}) {
//	s, ok := i.(string)
//	if !ok {
//		f.Errors(fmt.Sprint("%v is not a string", s))
//	}
//}

func FruitAndVegetable() Former {
	//	f1.Validaters(AcceptableFruit)
	//f1.Filters(Isstring)
	//	f2.Validaters(UnacceptableVegetable)
	//f2.Filters(Isstring)
	return newform(f1, f2, f3)
}

var FruitAndVegetableForm Former = FruitAndVegetable()

//func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
//	req, _ := http.NewRequest(method, path, nil)
//	w := httptest.NewRecorder()
//	r.ServeHTTP(w, req)
//	return w
//}

type TestHandler struct {
	handlers map[string]http.HandlerFunc
}

func (t TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := t.handlers[r.Method]
	if h != nil {
		h(w, r)
	}
}

func testserve() *TestHandler {
	return &TestHandler{handlers: make(map[string]http.HandlerFunc)}
}

func testformget(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	f := FruitAndVegetableForm.New()
	out := bytes.NewBuffer([]byte(`<form action="/" method="POST">`))
	out.WriteString(f.Render())
	out.WriteString("</form>")
	w.Write(out.Bytes())
}

func testformpost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	//fmt.Printf("%+v\n", r)
	f := FruitAndVegetableForm.New()
	f.Process(r)
	//fmt.Printf("%+v\n", f)
	//for _, fd := range f.Fields() {
	//	fmt.Printf("%+v\n", fd.Data())
	//}
	out := bytes.NewBuffer([]byte(`<form action="/" method="POST">`))
	out.WriteString(f.Render())
	out.WriteString("</form>")
	w.Write(out.Bytes())
}

func TestGetForm(t *testing.T) {
	ts := testserve()
	ts.handlers["GET"] = testformget
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	//fmt.Printf("%+v\n", w.Body)
}

func TestPostForm(t *testing.T) {
	ts := testserve()
	ts.handlers["POST"] = testformpost
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	v, _ := url.ParseQuery("fruit=pear,apple,blueberry&vegetable=spinach")
	//fmt.Printf("%+v\n", v)
	req.PostForm = v
	w := httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	fmt.Printf("%+v\n", w)
}
