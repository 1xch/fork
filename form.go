package fork

import (
	"bytes"
	"html/template"
	"net/http"
)

type Form interface {
	New() Form
	Former
}

type Former interface {
	Fields(...Field) []Field
	Process(*http.Request)
	Buffer() *bytes.Buffer
	String() string
	Render() template.HTML
	Valid() bool
	Errors() []string
}

func NewForm(fs ...Field) Form {
	return &form{fields: fs}
}

type form struct {
	fields []Field
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

func (f *form) Valid() bool {
	for _, fd := range f.Fields() {
		if !fd.Valid(fd) {
			return false
		}
	}
	return true
}

func (f *form) Errors() []string {
	var ret []string
	for _, fd := range f.Fields() {
		ret = append(ret, fd.Errors(fd)...)
	}
	return ret
}
