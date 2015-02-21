package fork

import (
	"fmt"
	"net/http"
	"strings"
)

var textwidget Widget = NewWidget(`<input type="text" name="{{ .Name }}" value="{{ .Get }}">`)

func TextField(name string) Field {
	return &textfield{
		name:      name,
		data:      NewValue(""),
		processor: NewProcessor(textwidget, nil, nil),
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

func (t *textfield) Name(name ...string) string {
	if len(name) > 0 {
		t.name = strings.Join(name, "-")
	}
	return t.name
}

func (t *textfield) Get() *Value {
	return t.data
}

func (t *textfield) Set(r *http.Request) {
	val := r.FormValue(t.Name())
	t.data = NewValue(val)
}

func TextAreaWidget(rows, cols int) Widget {
	ta := fmt.Sprintf(`<textarea name="{{ .Name }}" rows="%d" cols"%d">{{ .Get }}</textarea>`, rows, cols)
	return NewWidget(ta)
}

func TextAreaField(name string, rows, cols int) Field {
	return &textfield{
		name:      name,
		data:      NewValue(""),
		processor: NewProcessor(TextAreaWidget(rows, cols), nil, nil),
	}
}

var hiddenwidget Widget = NewWidget(`<input type="hidden" name="{{ .Name }}" value="{{ .Get }}">`)

func HiddenField(name string) Field {
	return &textfield{
		name:      name,
		data:      NewValue(""),
		processor: NewProcessor(hiddenwidget, nil, nil),
	}
}
