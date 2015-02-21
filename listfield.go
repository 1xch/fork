package fork

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func listfieldswidget(name string) Widget {
	lfw := fmt.Sprintf(`<fieldset name="%s"><ul>{{ range $x := .Get.Raw }}<li>{{ .Render $x }}</li>{{ end }}</ul></fieldset>`, name)
	return NewWidget(lfw)
}

func renamefield(name string, number int, field Field) Field {
	field.Name(name, strconv.Itoa(number), field.Name())
	return field
}

func renamefields(name string, number int, field Field) []Field {
	var ret []Field
	for i := 0; i < number; i++ {
		fd := field.New()
		renamefield(name, i, fd)
		ret = append(ret, fd)
	}
	return ret
}

func getregexp(name string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s-", name))
}

func ListField(name string, startwith int, start Field) Field {
	return &listfield{
		name:      name,
		base:      start,
		fields:    renamefields(name, startwith, start),
		processor: NewProcessor(listfieldswidget(name), nil, nil),
	}
}

type listfield struct {
	name   string
	base   Field
	fields []Field
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
	return NewValue(lf.fields)
}

func (lf *listfield) Set(r *http.Request) {
	lf.fields = nil
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
			lf.fields = append(lf.fields, nf)
		}
	}
}
