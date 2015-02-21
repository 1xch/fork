package fork

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func SimpleForm() Form {
	return NewForm(TextField("formfieldtext"), BoolField("formfieldbool", "FormFieldBool", false))
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

func PerformForForm(f Form, postdata string) (*httptest.ResponseRecorder, *httptest.ResponseRecorder) {
	ts := testserve()

	ts.handlers["GET"] = getformhandlerfor(f)
	ts.handlers["POST"] = postformhandlerfor(f)

	w1 := PerformGet(ts)
	w2 := PerformPost(ts, postdata)

	return w1, w2
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
	f := NewTestForm(TextField("text"))
	w1, w2 := PerformForForm(f, `text=TEXT`)
	fmt.Printf("%+v\n\n\n%+v\n\n\n", w1, w2)
}

func TestBoolField(t *testing.T) {
	f := NewTestForm(BoolField("yes", "YES", true), BoolField("no", "NO", false))
	w1, w2 := PerformForForm(f, `yes=false&no=true`)
	fmt.Printf("%+v\n\n\n%+v\n\n\n", w1, w2)
}

func TestRadioField(t *testing.T) {
	f := NewTestForm(RadioField("radiofield", "UP", "up", true), RadioField("radiofield", "DOWN", "down", false))
	w1, w2 := PerformForForm(f, ``)
	fmt.Printf("%+v\n\n\n%+v\n\n\n", w1, w2)
}

func TestCheckField(t *testing.T) {
	f := NewTestForm(CheckField("checkfield-1", "left"), CheckField("checkfield-2", "right"))
	w1, w2 := PerformForForm(f, ``)
	fmt.Printf("%+v\n\n\n%+v\n\n\n", w1, w2)
}

func TestSelectField(t *testing.T) {
	var testoptions []Option = []Option{
		[2]string{"one", "ONE"},
		[2]string{"two", "TWO"},
		[2]string{"three", "3"},
	}

	f := NewTestForm(SelectField("selectfield", testoptions...))

	w1, w2 := PerformForForm(f, ``)
	fmt.Printf("%+v\n\n\n%+v\n\n\n", w1, w2)
}

func TestListField(t *testing.T) {
	var ListField1 Field = ListField("listfield", 3, TextField("TEST"))

	f := NewTestForm(ListField1)

	w1, w2 := PerformForForm(f, `listfield-0-TEST=IamZERO&listfield-1-TEST=IamONE&listfield7-seven=IshouldnotbeSEVEN`)
	fmt.Printf("%+v\n\n\n%+v\n\n\n", w1, w2)
}

func TestFormsField(t *testing.T) {
	var FormsField1 Field = FormsField("formfield", 2, SimpleForm())

	f := NewTestForm(FormsField1)

	w1, w2 := PerformForForm(f, `formfield-0-formfieldtext-0=FORMFIELD0&formfield-0-formfieldbool-1=true&formfield-1-formfieldtext-0=FORMFIELD1&formfield-1-formfieldbool-1=false`)
	fmt.Printf("%+v\n\n\n%+v\n\n\n", w1, w2)
}
