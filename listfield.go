package fork

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func listfieldswidget(options ...string) Widget {
	return NewWidget(fmt.Sprintf(`<fieldset name="{{ .Name }}" %s><ul>{{ range $x := .Fields }}<li>{{ .Render $x }}</li>{{ end }}</ul></fieldset>`, strings.Join(options, " ")))
}

func renamefield(name string, number int, field Field) Field {
	field.Name(name, strconv.Itoa(number), field.Name())
	return field
}

func renamefields(name string, number int, field Field) []Field {
	var ret []Field
	for i := 0; i < number; i++ {
		fd := field.New()
		ret = append(ret, renamefield(name, i, fd))
	}
	return ret
}

func getregexp(name string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s-", name))
}

func ListField(name string, startwith int, start Field, options ...string) Field {
	return &listfield{
		name:      name,
		base:      start,
		Fields:    renamefields(name, startwith, start),
		processor: NewProcessor(listfieldswidget(options...), nil, nil),
	}
}

type listfield struct {
	name   string
	base   Field
	Fields []Field
	*processor
}

func (lf *listfield) New() Field {
	var newfield listfield = *lf
	return &newfield
}

func (lf *listfield) Name(name ...string) string {
	if len(name) > 0 {
		lf.name = strings.Join(name, "-")
	}
	return lf.name
}

func (lf *listfield) Get() *Value {
	return NewValue(lf.Fields)
}

func (lf *listfield) Set(r *http.Request) {
	lf.Fields = nil
	pl := 0
	for k, _ := range r.PostForm {
		if getregexp(lf.name).MatchString(k) {
			pl++
		}
	}
	if pl > 0 {
		for x := 0; x < pl; x++ {
			nf := lf.base.New()
			renamefield(lf.name, x, nf)
			nf.Set(r)
			lf.Fields = append(lf.Fields, nf)
		}
	}
}
