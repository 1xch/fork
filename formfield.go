package fork

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func formfieldwidget(name string) Widget {
	lfw := fmt.Sprintf(`<fieldset name="%s"><input type="hidden" name="%s" value="{{ .Index }}"><ul>{{ range $x := .Get.Raw }}<li>{{ .Render }}</li>{{ end }}</ul></fieldset>`, name, name)
	return NewWidget(lfw)
}

func renameformfields(name string, number int, f Form) Form {
	for index, field := range f.Fields() {
		field.Name(name, strconv.Itoa(number), field.Name(), strconv.Itoa(index))
	}
	return f
}

func addform(name string, number int, form Form) []Form {
	var ret []Form
	for i := 0; i < number; i++ {
		n := form.New()
		ret = append(ret, renameformfields(name, i, n))

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
		Index:     startwith,
		forms:     addform(name, startwith, start),
		processor: NewProcessor(formfieldwidget(name), nil, nil),
	}
}

type formfield struct {
	name  string
	base  Form
	Index int
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
	if len(r.PostForm) == 0 {
		r.ParseForm()
	}
	ff.forms = nil
	index := r.FormValue(ff.Name())
	i, err := strconv.Atoi(index)
	if err != nil {
		ff.Errors("form field index error: %s", err.Error())
	}
	ff.Index = i
	for x := 0; x < ff.Index; x++ {
		nf := ff.base.New()
		renameformfields(ff.name, x, nf)
		nf.Process(r)
		ff.forms = append(ff.forms, nf)
	}
}
