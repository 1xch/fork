package fork

import (
	"bytes"
	"net/http"
)

type Form interface {
	New() Form
	Former
}

type Former interface {
	Fields(...Field) []Field
	Process(*http.Request)
	Render() string
	Valid() bool
}

func NewForm(fs ...Field) *defaultform {
	return &defaultform{fields: fs}
}

type defaultform struct {
	//Validater
	fields []Field
}

func (df *defaultform) New() Form {
	var newform defaultform = *df
	fs := df.fields
	newform.fields = nil
	for _, field := range fs {
		newform.fields = append(newform.fields, field.New())
	}
	return &newform
}

func (df *defaultform) Fields(fs ...Field) []Field {
	df.fields = append(df.fields, fs...)
	return df.fields
}

func (df *defaultform) Process(r *http.Request) {
	for _, fd := range df.Fields() {
		fd.Set(r)
	}
}

func (df *defaultform) Valid() bool {
	for _, fd := range df.Fields() {
		if !fd.Valid() {
			return false
		}
	}
	return true
}

func (df *defaultform) Render() string {
	b := new(bytes.Buffer)
	for _, fd := range df.Fields() {
		_, _ = b.WriteString(fd.Render(fd))
	}
	return b.String()
}
