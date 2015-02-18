package fork

import "strconv"

var booleanwidget Widget = NewDefaultWidget(`<input type="checkbox" name="{{.Name}}"{{if .Get.Raw}} checked{{end}}>`)

func BooleanField(name string) *booleanfield {
	return &booleanfield{
		name:      name,
		checked:   NewValue(false),
		processor: DefaultProcessor(booleanwidget),
	}
}

type booleanfield struct {
	name    string
	checked *Value
	*processor
}

func (b *booleanfield) New() Field {
	var newfield booleanfield = *b
	return &newfield
}

func (b *booleanfield) Name() string {
	return b.name
}

func (b *booleanfield) Get() *Value {
	return b.checked
}

func (b *booleanfield) Set(i interface{}) {
	set, err := strconv.ParseBool(i.(string))
	if err != nil {
		b.checked = NewValue(false)
	}
	b.checked = NewValue(set)
}
