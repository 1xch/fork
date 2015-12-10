package fork

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var formhead = []byte(`<form action="/" method="POST">`)
var formtail = []byte(`</form>`)

func wrapForm(f Form) *bytes.Buffer {
	out := new(bytes.Buffer)
	out.Write(formhead)
	out.ReadFrom(f.Buffer())
	out.Write(formtail)
	return out
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

func PerformForm(t *testing.T, f Form, postData string) (*httptest.ResponseRecorder, *httptest.ResponseRecorder) {
	ts := testServe()

	ts.handlers["GET"] = formHandlerGet(t, f)
	ts.handlers["POST"] = formHandlerPost(t, f)

	for _, fd := range f.Fields() {
		if fd.Name() == "test" {
			tf := fd.Render(fd)
			if !strings.Contains(string(tf), "<input type=\"text\" name=\"test\" value=\"\" >") {
				t.Errorf("test field was: %s", tf)
			}
		}
	}

	w1 := PerformGet(ts)
	w2 := PerformPost(ts, postData)

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

func testServe() *TestHandler {
	return &TestHandler{handlers: make(map[string]http.HandlerFunc)}
}

func formHandlerGet(t *testing.T, f Form) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		nf := f.New()
		if f.Tag() != nf.Tag() {
			t.Errorf(
				"provided form and new form tags are not the same: %s %s",
				f.Tag(),
				nf.Tag(),
			)
		}
		out := wrapForm(nf)
		w.Write(out.Bytes())
	}
}

func formHandlerPost(t *testing.T, f Form) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		nf := f.New()
		w.WriteHeader(200)
		nf.Process(r)
		if vals := nf.Values(); vals == nil {
			t.Error("form.Values() did not return properly")
		}
		out := wrapForm(nf)
		w.Write(out.Bytes())
	}
}

func bodyExpectation(t *testing.T, r *httptest.ResponseRecorder, tag string, expectation string) {
	if !strings.Contains(r.Body.String(), expectation) {
		t.Errorf(
			"\n---\nform: %s\nhave:\n---\n%s\n\nexpected:\n---\n%s\n---\n",
			tag,
			r.Body,
			expectation,
		)
	}
}

type testField struct {
	Text      string
	validated bool
	*named
	Processor
}

func (f *testField) New(fc ...FieldConfig) Field {
	var newfield testField = *f
	newfield.named = f.named.Copy()
	newfield.Text = ""
	newfield.validated = false
	newfield.SetValidateable(false)
	return &newfield
}

func (f *testField) Get() *Value {
	return NewValue(f.Text)
}

func (f *testField) Set(r *http.Request) {
	v := f.Filter(f.Name(), r)
	f.Text = v.String()
	f.SetValidateable(true)
}

func testWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="text" name="{{ .Name }}" value="{{ .Text }}" %s>`, options...))
}

func FilterTest(in string) string {
	return "FILTERED"
}

func ValidaterTest(f *testField, t *testing.T) func(*testField) error {
	return func(f *testField) error {
		if f.Validateable() {
			if f.Text != "FILTERED" {
				t.Errorf("testField.Text should be FILTERED but was %s", f.Text)
			}
			f.validated = true
		}
		return nil
	}
}

func MakeTestField(t *testing.T, name string, options ...string) Field {
	tf := &testField{
		named: newnamed(name),
		Processor: NewProcessor(
			testWidget(options...),
			NewValidater(),
			NewFilterer(),
		),
	}
	tf.Filters(FilterTest)
	tf.Validaters(ValidaterTest(tf, t))
	return tf
}

func testBasic(t *testing.T, f Form, postProvides string, getExpects string, postExpects string) {
	f.Fields(MakeTestField(t, "test"))

	w1, w2 := PerformForm(t, f, postProvides)

	if w1.Code != 200 || w2.Code != 200 {
		t.Errorf(
			"Response incorrect; received Get %d Post %d, expected 200",
			w1.Code,
			w2.Code,
		)
	}

	bodyExpectation(t, w1, f.Tag(), getExpects)
	bodyExpectation(t, w2, f.Tag(), postExpects)
}
