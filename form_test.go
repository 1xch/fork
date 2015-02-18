package fork

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"
)

var (
	TestForm Form = NewTestForm()
)

func NewTestForm() Form {
	return &testform{
		fields: []Field{
			TextField("testtext"),
			TextField("email"),
			BooleanField("testbool"),
			HiddenField("invisible"),
			TextAreaField("textareatest", 10, 10),
		},
	}
}

type testform struct {
	//Validater
	fields []Field
}

func (f *testform) New() Form {
	var newform testform = *f
	//fs := f.fields
	//newform.fields = nil
	//for _, field := range fs {
	//	newform.fields = append(newform.fields, field.New())
	//}
	return &newform
}

func (f *testform) Fields(fs ...Field) []Field {
	f.fields = append(f.fields, fs...)
	return f.fields
}

func (f *testform) Process(r *http.Request) {
	for _, fd := range f.Fields() {
		//fmt.Printf("%+v\n", r.FormValue(fd.Name()))
		fd.Set(r.FormValue(fd.Name()))
	}
}

func (f *testform) Valid() bool {
	for _, fd := range f.Fields() {
		if !fd.Valid() {
			return false
		}
	}
	return true
}

func (f *testform) Render() string {
	b := new(bytes.Buffer)
	for _, fd := range f.Fields() {
		_, _ = b.WriteString(fd.Render(fd))
	}
	return b.String()
}

type TestHandler struct {
	handlers map[string]http.HandlerFunc
}

func (t TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := t.handlers[r.Method]
	if h != nil {
		h(w, r)
	}
}

func testserve() *TestHandler {
	return &TestHandler{handlers: make(map[string]http.HandlerFunc)}
}

func testformget(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	t := TestForm.New()
	out := bytes.NewBuffer([]byte(`<form action="/" method="POST">`))
	out.WriteString(t.Render())
	out.WriteString("</form>")
	w.Write(out.Bytes())
}

func testfieldrender(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(200)
	//t := TestForm.New()
	//out := bytes.NewBuffer([]byte(`<form action="/" method="POST">`))
	//out.WriteString(t.Render())
	//out.WriteString("</form>")
	//w.Write(out.Bytes())
}

func testformpost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	t := TestForm.New()
	t.Process(r)
	//for _, fd := range t.Fields() {
	//	fmt.Printf("%#v\n", fd.Get())
	//}
	out := bytes.NewBuffer([]byte(`<form action="/" method="POST">`))
	out.WriteString(t.Render())
	out.WriteString("</form>")
	w.Write(out.Bytes())
}

func TestGetForm(t *testing.T) {
	ts := testserve()
	ts.handlers["GET"] = testformget
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	//fmt.Printf("%+v\n", w.Body)
}

func TestPostForm(t *testing.T) {
	ts := testserve()
	ts.handlers["POST"] = testformpost
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	v, _ := url.ParseQuery("testtext=TESTTEXT1&email=test@test.com&testbool=true&invisible=unseen&textareatest='4scoreN7yearsago'")
	//fmt.Printf("%+v\n", v)
	req.PostForm = v
	w := httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	fmt.Printf("%+v\n", w.Body)
}
