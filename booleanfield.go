package fork

import (
	"net/http"
	"strconv"
	"strings"
)

type Option [2]string

type OptionInfo struct {
	Value string
	Label string
	Set   bool
}

func optioninfo(value string, label string, set bool) *OptionInfo {
	return &OptionInfo{value, label, set}
}

type booleanfield struct {
	name    string
	label   string
	data    *Value
	setfunc func(*booleanfield, *http.Request)
	*processor
}

func (b *booleanfield) New() Field {
	var newfield booleanfield = *b
	return &newfield
}

func (b *booleanfield) Name(name ...string) string {
	if len(name) > 0 {
		b.name = strings.Join(name, "-")
	}
	return b.name
}

func (b *booleanfield) Get() *Value {
	return b.data
}

func (b *booleanfield) Set(req *http.Request) {
	b.setfunc(b, req)
}

var boolwidget Widget = NewWidget(`<input type="checkbox" name="{{ .Name }}"{{ if .Get.Raw.Set }} checked{{ end }}>{{ .Get.Raw.Label }}`)

func BoolField(name string, label string, start bool) Field {
	return &booleanfield{
		name:      name,
		label:     label,
		data:      NewValue(optioninfo(name, label, start)),
		setfunc:   boolset,
		processor: NewProcessor(boolwidget, nil, nil),
	}
}

func boolset(b *booleanfield, req *http.Request) {
	val := req.FormValue(b.Name())
	set, err := strconv.ParseBool(val)
	if err != nil {
		set = false
	}
	b.data = NewValue(optioninfo(b.name, b.label, set))
}

var radiowidget Widget = NewWidget(`<input type="radio" name="{{ .Name }}" value="{{ .Get.Raw.Value }}"{{ if .Get.Raw.Set }} checked{{ end }}>{{ .Get.Raw.Label }}`)

func RadioField(name string, label string, value string, checked bool) Field {
	return &booleanfield{
		name:      name,
		label:     label,
		data:      NewValue(optioninfo(value, label, checked)),
		setfunc:   radioset,
		processor: NewProcessor(radiowidget, nil, nil),
	}
}

func radioset(b *booleanfield, req *http.Request) {}

var checkboxwidget Widget = NewWidget(`<input type="checkbox" name="{{ .Name }}" value="{{ .Get.Raw.Value }}">{{ .Get.Raw.Label }}`)

func CheckField(name string, label string) Field {
	return &booleanfield{
		name:      name,
		label:     label,
		data:      NewValue(optioninfo("", label, false)),
		setfunc:   checkset,
		processor: NewProcessor(checkboxwidget, nil, nil),
	}
}

func checkset(b *booleanfield, req *http.Request) {}
