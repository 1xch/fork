package fork

import "fmt"

var textwidget Widget = NewDefaultWidget(`<input type="text" name="{{ .Name }}" value="{{ .Get }}">`)

func TextField(name string) *textfield {
	return &textfield{
		name:      name,
		data:      NewValue(""),
		processor: DefaultProcessor(textwidget),
	}
}

type textfield struct {
	name string
	data *Value
	*processor
}

func (t *textfield) New() Field {
	var newfield textfield = *t
	return &newfield
}

func (t *textfield) Name() string {
	return t.name
}

func (t *textfield) Get() *Value {
	return t.data
}

func (t *textfield) Set(i interface{}) {
	t.data = NewValue(i)
}

func TextAreaWidget(rows, cols int) Widget {
	ta := fmt.Sprintf(`<textarea name="{{ .Name }}" rows="%d" cols"%d">{{ .Get }}</textarea>`, rows, cols)
	return NewDefaultWidget(ta)
}

func TextAreaField(name string, rows, cols int) *textfield {
	return &textfield{
		name:      name,
		data:      NewValue(""),
		processor: DefaultProcessor(TextAreaWidget(rows, cols)),
	}
}

var hiddenwidget Widget = NewDefaultWidget(`<input type="hidden" name="{{ .Name }}" value="{{ .Get }}">`)

func HiddenField(name string) *textfield {
	return &textfield{
		name:      name,
		data:      NewValue(""),
		processor: DefaultProcessor(hiddenwidget),
	}
}
