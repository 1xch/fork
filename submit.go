package fork

import "net/http"

type submitField struct {
	Submitted bool
	*baseField
	*processor
}

func (s *submitField) New() Field {
	var newfield submitField = *s
	newfield.baseField = s.baseField.Copy()
	newfield.Submitted = false
	newfield.validateable = false
	return &newfield
}

func (s *submitField) Get() *Value {
	return NewValue(s.Submitted)
}

func (s *submitField) Set(r *http.Request) {
	s.Filter(s.Name(), r)
	s.Submitted = true
	s.validateable = true
}

func submitWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="submit" name="{{ .Name }}" value="{{ .Name }}" %s>`, options...))
}

func SubmitField(name string, v []interface{}, f []interface{}, options ...string) Field {
	return &submitField{
		baseField: newBaseField(name),
		processor: NewProcessor(
			submitWidget(options...),
			NewValidater(v...),
			NewFilterer(f...),
		),
	}
}
