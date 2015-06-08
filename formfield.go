package fork

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func NewFieldIndex(s string, i int) *FieldIndex {
	return &FieldIndex{N: s, I: i}
}

type FieldIndex struct {
	I int
	N string
}

type formField struct {
	base  Form
	Index *FieldIndex
	Forms []Form
	*baseField
	*processor
}

func (f *formField) New() Field {
	var newfield formField = *f
	newfield.baseField = f.baseField.Copy()
	newfield.validateable = false
	return &newfield
}

func (f *formField) Get() *Value {
	return NewValue(f.Forms)
}

func (f *formField) Set(r *http.Request) {
	f.Forms = nil
	i := f.Filter(f.Name(), r)
	f.Index = i.Raw.(*FieldIndex)
	for x := 0; x < f.Index.I; x++ {
		nf := f.base.New()
		renameFormFields(f.name, x, nf)
		nf.Process(r)
		f.Forms = append(f.Forms, nf)
	}
	f.validateable = true
}

func FilterIndex(index string) *FieldIndex {
	i, _ := strconv.Atoi(index)
	return NewFieldIndex(index, i)
}

func ValidateIndex(f *formField) error {
	_, err := strconv.Atoi(f.Index.N)
	if err != nil {
		return fmt.Errorf("form field index error: %s", err.Error())
	}
	return nil
}

func formFieldWidget(name string) Widget {
	in := strings.Join([]string{
		fmt.Sprintf(`<fieldset name="%s">`, name),
		fmt.Sprintf(`<input type="hidden" name="%s" `, name),
		`value="{{ .Index.N }}"><ul>{{ range $x := .Forms }}`,
		`<li>{{ .Render }}</li>{{ end }}</ul></fieldset>`,
	}, "")
	return NewWidget(in)
}

func renameFormFields(name string, number int, f Form) Form {
	for index, field := range f.Fields() {
		field.Name(name, strconv.Itoa(number), field.Name(), strconv.Itoa(index))
	}
	return f
}

func addForm(name string, number int, form Form) []Form {
	var ret []Form
	for i := 0; i < number; i++ {
		n := form.New()
		ret = append(ret, renameFormFields(name, i, n))

	}
	return ret
}

func FormField(name string, f Form) Field {
	return FormsField(name, 1, f)
}

func FormsField(name string, startwith int, start Form) Field {
	return &formField{
		base:      start,
		Index:     NewFieldIndex(strconv.Itoa(startwith), startwith),
		Forms:     addForm(name, startwith, start),
		baseField: newBaseField(name),
		processor: NewProcessor(
			formFieldWidget(name),
			NewValidater(ValidateIndex),
			NewFilterer(FilterIndex),
		),
	}
}
