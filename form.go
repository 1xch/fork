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
	String() string
	Render() template.HTML
	Validater
}

func NewForm(fs ...Field) Form {
	return &form{fields: fs}
}

type form struct {
	fields []Field
	Validater
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
	for _, fd := range f.Fields() {
		fd.Set(r)
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

func (f *form) Validate(fd Field) error {
	return fd.Validate(fd)
}

func (f *form) String() string {
	b := new(bytes.Buffer)
	for _, fd := range f.Fields() {
		_, _ = b.WriteString(fd.String(fd))
	}
	return b.String()
}

func (f *form) Render() template.HTML {
	return template.HTML(f.String())
}
