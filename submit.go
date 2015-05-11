package fork

import (
	"fmt"
	"net/http"
	"strings"
)

func SubmitField(name string, validaters []interface{}, filters []interface{}, options ...string) Field {
	return &submitfield{
		name: name,
		processor: NewProcessor(
			submitwidget(options...),
			validaters,
			filters,
		),
	}
}

func submitwidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<input type="submit" value="{{ .Name }}" %s>`, strings.Join(options, " ")))
}

type submitfield struct {
	name         string
	validateable bool
	*processor
}

func (s *submitfield) New() Field {
	var newfield submitfield = *s
	s.validateable = false
	return &newfield
}

func (s *submitfield) Name(name ...string) string {
	if len(name) > 0 {
		s.name = strings.Join(name, "-")
	}
	return s.name
}

func (s *submitfield) Get() *Value {
	return NewValue(nil)
}

func (s *submitfield) Set(r *http.Request) {
	s.validateable = true
}

func (s *submitfield) Validateable() bool {
	return s.validateable
}
