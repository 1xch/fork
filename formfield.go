package fork

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func formfieldwidget(name string) Widget {
	lfw := fmt.Sprintf(`<fieldset name="%s"><input type="hidden" name="%s" value="{{ .Index.N }}"><ul>{{ range $x := .Forms }}<li>{{ .Render }}</li>{{ end }}</ul></fieldset>`, name, name)
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
		name:  name,
		base:  start,
		Index: NewFieldIndex(strconv.Itoa(startwith), startwith),
		Forms: addform(name, startwith, start),
		processor: NewProcessor(
			formfieldwidget(name),
			[]interface{}{ValidateIndex},
			[]interface{}{FilterIndex},
		),
	}
}

func NewFieldIndex(s string, i int) *FieldIndex {
	return &FieldIndex{N: s, I: i}
}

type FieldIndex struct {
	I int
	N string
}

type formfield struct {
	name         string
	base         Form
	Index        *FieldIndex
	Forms        []Form
	validateable bool
	*processor
}

func (ff *formfield) New() Field {
	var newfield formfield = *ff
	ff.validateable = false
	return &newfield
}

func (ff *formfield) Name(name ...string) string {
	if len(name) > 0 {
		ff.name = strings.Join(name, "-")
	}
	return ff.name
}

func (ff *formfield) Get() *Value {
	return NewValue(ff.Forms)
}

func (ff *formfield) Set(r *http.Request) {
	ff.Forms = nil
	i := ff.Filter(ff.Name(), r)
	ff.Index = i.Raw.(*FieldIndex)
	for x := 0; x < ff.Index.I; x++ {
		nf := ff.base.New()
		renameformfields(ff.name, x, nf)
		nf.Process(r)
		ff.Forms = append(ff.Forms, nf)
	}
	ff.validateable = true
}

func (ff *formfield) Validateable() bool {
	return ff.validateable
}

func FilterIndex(index string) *FieldIndex {
	i, _ := strconv.Atoi(index)
	return NewFieldIndex(index, i)
}

func ValidateIndex(ff *formfield) error {
	_, err := strconv.Atoi(ff.Index.N)
	if err != nil {
		return fmt.Errorf("form field index error: %s", err.Error())
	}
	return nil
}
