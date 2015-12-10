package fork

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type Widget interface {
	ToBytes(interface{}) (*bytes.Buffer, error)
	ToString(interface{}) string
	Render(Field) template.HTML
}

var defaultTemplate = strings.Join([]string{
	`{{ define "fielderrors" }}`,
	`<div class="field-errors">`,
	`<ul>{{ range $x := .Errors . }}`,
	`<li>{{ $x }}</li>`,
	`{{ end }}</ul></div>{{ end }}`,
	`{{ define "default" }}`,
	`%s`,
	`{{ if .Error .}}{{ template "fielderrors" .}}{{ end }}{{ end }}`,
}, "")

func WithOptions(base string, options ...string) string {
	return fmt.Sprintf(base, strings.Join(options, " "))
}

func NewWidget(t string) Widget {
	var err error
	ti := &widget{name: "default"}
	tt := fmt.Sprintf(defaultTemplate, t)
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

func (w *widget) ToBytes(i interface{}) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	err := w.widget.ExecuteTemplate(&buffer, w.name, i)
	if err != nil {
		return nil, err
	}
	return &buffer, nil
}

func (w *widget) ToString(i interface{}) string {
	b, err := w.ToBytes(i)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

func (w *widget) Render(f Field) template.HTML {
	return template.HTML(w.ToString(f))
}
