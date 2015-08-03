package fork

import "net/http"

type submitField struct {
	Submitted bool
	*named
	Processor
}

func (s *submitField) New(i ...interface{}) Field {
	var newfield submitField = *s
	newfield.named = s.named.Copy()
	newfield.Submitted = false
	newfield.SetValidateable(false)
	return &newfield
}

func (s *submitField) Get() *Value {
	return NewValue(s.Submitted)
}

func (s *submitField) Set(r *http.Request) {
	s.Filter(s.Name(), r)
	s.Submitted = true
	s.SetValidateable(true)
}

func submitWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="submit" name="{{ .Name }}" value="{{ .Name }}" %s>`, options...))
}

func SubmitField(name string, v []interface{}, f []interface{}, options ...string) Field {
	return &submitField{
		named: newnamed(name),
		Processor: NewProcessor(
			submitWidget(options...),
			NewValidater(v...),
			NewFilterer(f...),
		),
	}
}
