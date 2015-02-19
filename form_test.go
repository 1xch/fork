package fork

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func makenewtextfield(name string, req *http.Request) Field {
	t := TextField(name)
	t.Set(req)
	return t
}

func SimpleForm() Form {
	return NewForm(TextField("formfieldtext"), BooleanField("formfieldbool", "FFbool", false))
}

func NewTestForm(fds ...Field) Form {
	return NewForm(fds...)
}

func PerformGet(th *TestHandler) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	th.ServeHTTP(w, req)
	return w
}

func PerformPost(th *TestHandler, values string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	v, _ := url.ParseQuery(values)
	req.PostForm = v
	w := httptest.NewRecorder()
	th.ServeHTTP(w, req)
	return w
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

func getformhandlerfor(f Form) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		t := f.New()
		out := bytes.NewBuffer([]byte(`test page delivered by GET`))
		out.WriteString(`<form action="/" method="POST">`)
		out.WriteString(t.Render())
		out.WriteString("</form>")
		out.WriteString("end test page")
		w.Write(out.Bytes())
	}
}

func postformhandlerfor(f Form) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		t := f.New()
		t.Process(r)
		out := bytes.NewBuffer([]byte(`<form action="/" method="POST">`))
		out.WriteString(t.Render())
		out.WriteString("</form>")
		w.Write(out.Bytes())
	}
}

func TestTextField(t *testing.T) {
	ts := testserve()

	f := NewTestForm(TextField("text"))

	ts.handlers["GET"] = getformhandlerfor(f)
	ts.handlers["POST"] = postformhandlerfor(f)

	w := PerformGet(ts)
	fmt.Printf("%+v\n\n\n", w.Body)

	w = PerformPost(ts, `text=TEXT`)
	fmt.Printf("%+v\n\n\n", w.Body)
}

func TestBoolField(t *testing.T) {
	ts := testserve()

	f := NewTestForm(BooleanField("yes", "YES", true), BooleanField("no", "NO", false))

	ts.handlers["GET"] = getformhandlerfor(f)
	ts.handlers["POST"] = postformhandlerfor(f)

	w := PerformGet(ts)
	fmt.Printf("%+v\n\n\n", w.Body)

	w = PerformPost(ts, `yes=false&no=true`)
	fmt.Printf("%+v\n\n\n", w.Body)

}

var ListField1 Field = ListField("listfield", makenewtextfield, TextField("zero"))

func TestListField(t *testing.T) {
	ts := testserve()

	f := NewTestForm(ListField1)

	ts.handlers["GET"] = getformhandlerfor(f)
	ts.handlers["POST"] = postformhandlerfor(f)

	w := PerformGet(ts)
	fmt.Printf("%+v\n\n\n", w.Body)

	w = PerformPost(ts, `listfield-0-zero=IamZERO&listfield-1-one=IamONE&listfield7-seven=IshouldnotbeSEVEN`)
	fmt.Printf("%+v\n\n\n", w.Body)
}

var FormsField1 Field = FormsField("formfield", 2, SimpleForm())

func TestFormsField(t *testing.T) {
	ts := testserve()

	f := NewTestForm(FormsField1)

	ts.handlers["GET"] = getformhandlerfor(f)
	ts.handlers["POST"] = postformhandlerfor(f)

	w := PerformGet(ts)
	fmt.Printf("%+v\n\n\n", w.Body)

	w = PerformPost(ts, `formfield-0-formfieldtext-0=FORMFIELD0&formfield-0-formfieldbool-1=true&formfield-1-formfieldtext-0=FORMFIELD1&formfield-1-formfieldbool-1=false`)
	fmt.Printf("%+v\n\n\n", w.Body)
}
