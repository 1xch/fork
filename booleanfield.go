package fork

import (
	"net/http"
	"strconv"
	"strings"
)

var booleanwidget Widget = NewDefaultWidget(`<input type="checkbox" name="{{ .Name }}"{{ if .Get.Raw.Set }} checked{{ end }}>{{ .Get.Raw.Label }}`)

func BooleanField(name string, label string, start bool) *booleanfield {
	return &booleanfield{
		name:      name,
		label:     label,
		data:      NewValue(boolinfo(name, label, start)),
		processor: DefaultProcessor(booleanwidget),
	}
}

type booleanfield struct {
	name  string
	label string
	data  *Value
	*processor
}

type BoolInfo struct {
	Name  string
	Label string
	Set   bool
}

func boolinfo(name string, label string, value bool) *BoolInfo {
	return &BoolInfo{name, label, value}
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
	val := req.FormValue(b.Name())
	set, err := strconv.ParseBool(val)
	if err != nil {
		set = false //b.data = NewValue(boolinfo(b.name, b.label, false))
	}
	b.data = NewValue(boolinfo(b.name, b.label, set))
}
