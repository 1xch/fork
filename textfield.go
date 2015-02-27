package fork

import (
	"fmt"
	"net/http"
	"net/mail"
	"strings"
)

func textwidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<input type="text" name="{{ .Name }}" value="{{ .Text }}" %s>`, strings.Join(options, " ")))
}

func newtextfield(name string, widget Widget) Field {
	return &textfield{
		name:      name,
		Text:      "",
		processor: NewProcessor(widget, nil, nil),
	}
}

func TextField(name string, options ...string) Field {
	return newtextfield(name, textwidget(options...))
}

type textfield struct {
	name string
	Text string
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
	return NewValue(t.Text)
}

func (t *textfield) Set(r *http.Request) {
	val := r.FormValue(t.Name())
	t.Text = val
	err := t.Validate(t)
	if err != nil {
		t.Errors(err.Error())
	}
}

func textareawidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<textarea name="{{ .Name }}" %s>{{ .Text }}</textarea>`, strings.Join(options, " ")))
}

func TextAreaField(name string, options ...string) Field {
	return newtextfield(name, textareawidget(options...))
}

func hiddenwidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<input type="hidden" name="{{ .Name }}" value="{{ .Text }}" %s>`, strings.Join(options, " ")))
}

func HiddenField(name string, options ...string) Field {
	return newtextfield(name, hiddenwidget(options...))
}

func passwordwidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<input type="password" name="{{ .Name }}" value="{{ .Text }}" %s>`, strings.Join(options, " ")))
}

func PassWordField(name string, options ...string) Field {
	return newtextfield(name, passwordwidget(options...))
}

func emailwidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<input type="email" name="{{ .Name }}" value="{{ .Text }}" %s>`, strings.Join(options, " ")))
}

func EmailField(name string, options ...string) Field {
	return &textfield{
		name: name,
		Text: "",
		processor: NewProcessor(
			emailwidget(options...),
			[]interface{}{ValidateEmail},
			nil,
		),
	}
}

func ValidateEmail(f Field) error {
	a := f.Get()
	_, err := mail.ParseAddress(a.String())
	if err != nil {
		return fmt.Errorf("Invalid email address: %s", err.Error())
	}
	return nil
}
