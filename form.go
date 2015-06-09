package fork

import (
	"bytes"
	"html/template"
	"net/http"
)

type FormConfig func(Form) error

type Form interface {
	Tagged
	Creator
	Former
	Renderer
	Informer
	Checker
}

type Tagged interface {
	Tag() string
}

type Creator interface {
	New() Form
	Configure(...FormConfig) error
}

type Former interface {
	Fields(...Field) []Field
	Process(*http.Request)
}

type Renderer interface {
	Buffer() *bytes.Buffer
	String() string
	Render() template.HTML
}

type Informer interface {
	Values() map[string]*Value
	Errors() []string
}

func Checks(cs ...interface{}) FormConfig {
	return func(f Form) error {
		f.Checks(cs...)
		return nil
	}
}

func Fields(fs ...Field) FormConfig {
	return func(f Form) error {
		f.Fields(fs...)
		return nil
	}
}

func NewForm(tag string, fc ...FormConfig) Form {
	nf := &form{
		tag:     tag,
		Checker: NewChecker(),
	}
	nf.Configure(fc...)
	return nf
}

type form struct {
	tag    string
	fields []Field
	values map[string]*Value
	Checker
}

func (f *form) Tag() string {
	return f.tag
}

func (f *form) New() Form {
	var newform form = *f
	fs := f.fields
	newform.fields = nil
	for _, field := range fs {
		newform.fields = append(newform.fields, field.New())
	}
	return &newform
}

func (f *form) Configure(fc ...FormConfig) error {
	var err error
	for _, c := range fc {
		err = c(f)
	}
	return err
}

func (f *form) Fields(fs ...Field) []Field {
	f.fields = append(f.fields, fs...)
	return f.fields
}

func (f *form) Process(r *http.Request) {
	if len(r.PostForm) == 0 {
		r.ParseForm()
	}
	for _, fd := range f.Fields() {
		fd.Set(r)
	}
}

func (f *form) Buffer() *bytes.Buffer {
	b := new(bytes.Buffer)
	for _, fd := range f.Fields() {
		fb, err := fd.Bytes(fd)
		if err == nil {
			b.ReadFrom(fb)
		}
	}
	return b
}

func (f *form) String() string {
	return f.Buffer().String()
}

func (f *form) Render() template.HTML {
	return template.HTML(f.String())
}

func (f *form) mkvalues() {
	f.values = make(map[string]*Value)
	for _, fd := range f.fields {
		f.values[fd.Name()] = fd.Get()
	}
}

func (f *form) Values() map[string]*Value {
	if f.values == nil {
		f.mkvalues()
	}
	return f.values
}

func (f *form) Errors() []string {
	var ret []string
	for _, fd := range f.Fields() {
		ret = append(ret, fd.Errors(fd)...)
	}
	return ret
}
