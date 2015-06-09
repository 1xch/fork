package fork

import (
	"fmt"
	"net/http"
	"net/mail"
)

func newtextField(n *named, w Widget, v Validater, f Filterer) Field {
	return &textField{
		named:     n,
		Processor: NewProcessor(w, v, f),
	}
}

type textField struct {
	Text string
	*named
	Processor
}

func (t *textField) New() Field {
	var newfield textField = *t
	newfield.named = t.named.Copy()
	newfield.Text = ""
	newfield.SetValidateable(false)
	return &newfield
}

func (t *textField) Get() *Value {
	return NewValue(t.Text)
}

func (t *textField) Set(r *http.Request) {
	v := t.Filter(t.Name(), r)
	t.Text = v.String()
	t.SetValidateable(true)
}

func textWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="text" name="{{ .Name }}" value="{{ .Text }}" %s>`, options...))
}

func TextField(name string, v []interface{}, f []interface{}, options ...string) Field {
	return newtextField(
		newnamed(name),
		textWidget(options...),
		NewValidater(v...),
		NewFilterer(f...),
	)
}

func textAreaWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<textarea name="{{ .Name }}" %s>{{ .Text }}</textarea>`, options...))
}

func TextAreaField(name string, v []interface{}, f []interface{}, options ...string) Field {
	return newtextField(
		newnamed(name),
		textAreaWidget(options...),
		NewValidater(v...),
		NewFilterer(f...),
	)
}

func hiddenWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="hidden" name="{{ .Name }}" value="{{ .Text }}" %s>`, options...))
}

func HiddenField(name string, v []interface{}, f []interface{}, options ...string) Field {
	return newtextField(
		newnamed(name),
		hiddenWidget(options...),
		NewValidater(v...),
		NewFilterer(f...),
	)
}

func passwordWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="password" name="{{ .Name }}" value="{{ .Text }}" %s>`, options...))
}

func PassWordField(name string, v []interface{}, f []interface{}, options ...string) Field {
	return newtextField(
		newnamed(name),
		passwordWidget(options...),
		NewValidater(v...),
		NewFilterer(f...),
	)
}

func emailWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="email" name="{{ .Name }}" value="{{ .Text }}" %s>`, options...))
}

func EmailField(name string, v []interface{}, f []interface{}, options ...string) Field {
	return newtextField(
		newnamed(name),
		emailWidget(options...),
		NewValidater(append(v, ValidEmail)...),
		NewFilterer(f...),
	)
}

func ValidEmail(t *textField) error {
	if t.Validateable() {
		_, err := mail.ParseAddress(t.Text)
		if err != nil {
			return fmt.Errorf("Invalid email address: %s", err.Error())
		}
	}
	return nil
}
