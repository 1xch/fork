package fork

import (
	"fmt"
	"net/http"
	"strings"
)

func selectwidget(name string) Widget {
	ow := fmt.Sprintf(`<select name="%s">{{  range .Get.Raw }}<option value={{ .Value }}{{ if .Set }} selected{{ end }}>{{ .Label }}</option>{{ end }}</select>`, name)
	return NewWidget(ow)
}

func SelectOptions(selected Option, opts []Option) *Value {
	var val []*OptionInfo
	val = append(val, &OptionInfo{selected[0], selected[1], true})
	for _, o := range opts {
		val = append(val, &OptionInfo{o[0], o[1], false})
	}
	return NewValue(val)
}

func SelectField(name string, opts ...Option) Field {
	return &selectfield{
		name:      name,
		data:      SelectOptions(opts[0], opts[1:]),
		processor: NewProcessor(selectwidget(name), nil, nil),
	}
}

type selectfield struct {
	name string
	data *Value
	*processor
}

func (s *selectfield) New() Field {
	var newfield selectfield = *s
	return &newfield
}

func (s *selectfield) Name(name ...string) string {
	if len(name) > 0 {
		s.name = strings.Join(name, "-")
	}
	return s.name
}

func (s *selectfield) Get() *Value {
	return s.data
}

func (s *selectfield) Set(req *http.Request) {
	//val := req.FormValue(b.Name())
	//set, err := strconv.ParseBool(val)
	//if err != nil {
	//	set = false //b.data = NewValue(boolinfo(b.name, b.label, false))
	//}
	//b.data = NewValue(boolinfo(b.name, b.label, set))
}
