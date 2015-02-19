package fork

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func listfieldswidget(name string) Widget {
	lfw := fmt.Sprintf(`<fieldset name="%s"><ul>{{ range $x := .Get.Raw }}<li>{{ .Render $x }}</li>{{ end }}</ul></fieldset>`, name)
	return NewDefaultWidget(lfw)
}

type NewFieldFunc func(string, *http.Request) Field

func renamefields(name string, fields []Field) []Field {
	for index, field := range fields {
		field.Name(name, strconv.Itoa(index), field.Name())
	}
	return fields
}

func getregexp(name string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s-", name))
}

func ListField(name string, newfield NewFieldFunc, startingfield ...Field) *listfield {
	return &listfield{
		name:      name,
		newfield:  newfield,
		match:     getregexp(name),
		fields:    renamefields(name, startingfield),
		processor: DefaultProcessor(listfieldswidget(name)),
	}
}

type listfield struct {
	name     string
	newfield NewFieldFunc
	match    *regexp.Regexp
	fields   []Field
	*processor
}

func (lf *listfield) New() Field {
	var newfield listfield = *lf
	return &newfield
}

func (lf *listfield) Name(name ...string) string {
	if len(name) > 0 {
		lf.name = strings.Join(name, "-")
	}
	return lf.name
}

func (lf *listfield) Get() *Value {
	return NewValue(lf.fields)
}

func (lf *listfield) Set(r *http.Request) {
	for k, _ := range r.PostForm {
		if lf.match.MatchString(k) {
			n := lf.newfield(k, r)
			lf.fields = append(lf.fields, n)
		}
	}
}
