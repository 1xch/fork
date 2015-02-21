package fork

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func formfieldwidget(name string) Widget {
	lfw := fmt.Sprintf(`<fieldset name="%s"><ul>{{ range $x := .Get.Raw }}<li>{{ .Render }}</li>{{ end }}</ul></fieldset>`, name)
	return NewWidget(lfw)
}

func renameformfield(name string, number int, f Form) {
	for index, field := range f.Fields() {
		field.Name(name, strconv.Itoa(number), field.Name(), strconv.Itoa(index))
	}
}

func addform(name string, number int, form Form) []Form {
	var ret []Form
	for i := 0; i < number; i++ {
		n := form.New()
		renameformfield(name, i, n)
		ret = append(ret, n)

	}
	return ret
}

func FormField(name string, f Form) Field {
	return FormsField(name, 1, f)
}

func FormsField(name string, startwith int, start Form) Field {
	return &formfield{
		name:      name,
		base:      start,
		forms:     addform(name, startwith, start),
		processor: NewProcessor(formfieldwidget(name), nil, nil),
	}
}

type formfield struct {
	name  string
	base  Form
	forms []Form
	*processor
}

func (ff *formfield) New() Field {
	var newfield formfield = *ff
	return &newfield
}

func (ff *formfield) Name(name ...string) string {
	if len(name) > 0 {
		ff.name = strings.Join(name, "-")
	}
	return ff.name
}

func (ff *formfield) Get() *Value {
	return NewValue(ff.forms)
}

func (ff *formfield) Set(r *http.Request) {
	ff.forms = nil
	pl := 0
	for k, _ := range r.PostForm {
		if getregexp(ff.name).MatchString(k) {
			pl++
		}
	}
	if pl > 0 {
		al := pl / len(ff.base.Fields())
		for x := 0; x < al; x++ {
			nf := ff.base.New()
			renameformfield(ff.name, x, nf)
			nf.Process(r)
			ff.forms = append(ff.forms, nf)
		}
	}
}
