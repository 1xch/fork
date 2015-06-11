package fork

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

var formhead = []byte(`<form action="/" method="POST">`)
var formtail = []byte(`</form>`)

func formbuffer(note string, f Form) *bytes.Buffer {
	out := bytes.NewBuffer([]byte(fmt.Sprintf("%s\n", note)))
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

func testCheck(t *testing.T) func(Form) (bool, error) {
	return func(f Form) (bool, error) {
		for _, fd := range f.Fields() {
			if !fd.Valid(fd) {
				t.Errorf("invalid field %+v in test check", fd)
				return false, errors.New("invalid field")
			}
		}
		return true, nil
	}
}

func PerformForForm(t *testing.T, f Form, postdata string) (*httptest.ResponseRecorder, *httptest.ResponseRecorder) {
	ts := testserve()

	f.Checks(testCheck(t))

	ts.handlers["GET"] = formHandlerGet(t, f)
	ts.handlers["POST"] = formHandlerPost(t, f)

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

func formHandlerGet(t *testing.T, f Form) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		nf := f.New()
		if f.Tag() != nf.Tag() {
			t.Errorf("provided form and new form tags are not the same: %s %s", f.Tag(), nf.Tag())
		}
		out := formbuffer(`form via GET`, nf)
		w.Write(out.Bytes())
	}
}

func formHandlerPost(t *testing.T, f Form) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		nf := f.New()
		if ok, err := nf.Check(nf); ok && err == nil {
			w.WriteHeader(200)
			nf.Process(r)
			if vals := nf.Values(); vals == nil {
				t.Error("form.Values() did not return properly")
			}
			out := formbuffer(`form via POST`, nf)
			w.Write(out.Bytes())
		} else {
			t.Errorf("form %+v is not ok and/or contains errors: %+v", nf, err)
		}
	}
}

func testBasic(t *testing.T, f Form, postprovides string, GETexpects string, POSTexpects string) {
	w1, w2 := PerformForForm(t, f, postprovides)

	if w1.Code != 200 || w2.Code != 200 {
		t.Errorf("Response incorrect; received Get %d Post %d, expected 200", w1.Code, w2.Code)
	}

	if !strings.Contains(w1.Body.String(), GETexpects) {
		t.Errorf("\n---\n%s GET Error\nhave:\n---\n%s\n\nexpected:\n---\n%s\n---\n", f.Tag(), w1.Body, GETexpects)
	}

	if !strings.Contains(w2.Body.String(), POSTexpects) {
		t.Errorf("\n---\n%s POST Error\nhave:\n---\n%s\n\nexpected:\n---\n%s\n---\n", f.Tag(), w2.Body, POSTexpects)
	}

}

func testvaluestring(t *testing.T, expected, provided string) {
	if provided != expected {
		t.Errorf("Expected %s, but received %s from form.Values() value.String()", expected, provided)
	}
}

func testvaluebool(t *testing.T, expected, provided bool) {
	if provided != expected {
		t.Errorf("Expected %s, but received %s from form.Values() value.Bool()", expected, provided)
	}
}

func TestValues(t *testing.T) {
	f := NewForm(
		"TestValues",
		Fields(
			TextField("ValueOne", nil, nil),
			TextField("ValueTwo", nil, nil),
			BooleanField("yes", "YES", true),
		),
	)
	v1, v2, v3 := "VALUE1", "VALUE2", "no"
	post := fmt.Sprintf(`ValueOne=%s&ValueTwo=%s&yes=%s`, v1, v2, v3)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		ff := f.New()
		ff.Process(r)
		v := ff.Values()
		testvaluestring(t, "VALUE1", v["ValueOne"].String())
		testvaluestring(t, "VALUE2", v["ValueTwo"].String())
		testvaluebool(t, false, v["yes"].Bool())
		out := formbuffer(`form via POST`, ff)
		w.Write(out.Bytes())
	}
	ts := testserve()
	ts.handlers["POST"] = handler
	PerformPost(ts, post)
}

func TestTextField(t *testing.T) {
	testBasic(
		t,
		NewForm("TextField", Fields(TextField("text", nil, nil))),
		`text=TEXT`,
		`<input type="text" name="text" value="" >`,
		`<input type="text" name="text" value="TEXT" >`,
	)
}

func TestTextAreaField(t *testing.T) {
	testBasic(
		t,
		NewForm("TextAreaField", Fields(TextAreaField("textarea", nil, nil, "rows=10", "cols=10"))),
		`textarea=TEXTAREA`,
		`<textarea name="textarea" rows=10 cols=10></textarea>`,
		`<textarea name="textarea" rows=10 cols=10>TEXTAREA</textarea>`,
	)
}

func TestHiddenField(t *testing.T) {
	testBasic(
		t,
		NewForm("HiddenField", Fields(HiddenField("hidden", nil, nil))),
		`hidden=HIDDEN`,
		`<input type="hidden" name="hidden" value="" >`,
		`<input type="hidden" name="hidden" value="HIDDEN" >`,
	)
}

func TestPasswordField(t *testing.T) {
	testBasic(
		t,
		NewForm("PasswordField", Fields(PassWordField("password", nil, nil, "size=10", "maxlength=30"))),
		`password=PASSWORD`,
		`<input type="password" name="password" value="" size=10 maxlength=30>`,
		`<input type="password" name="password" value="PASSWORD" size=10 maxlength=30>`,
	)
}

func TestEmailField(t *testing.T) {
	testBasic(
		t,
		NewForm("EmailField", Fields(EmailField("email", nil, nil))),
		`email=test@test.com`,
		`<input type="email" name="email" value="" >`,
		`<input type="email" name="email" value="test@test.com" >`,
	)
	testBasic(
		t,
		NewForm("EmailField :: Error", Fields(EmailField("email", nil, nil))),
		`email=invalidemail.com`,
		`<input type="email" name="email" value="" >`,
		`<input type="email" name="email" value="invalidemail.com" ><div class="field-errors"><ul><li>Invalid email address: mail: missing phrase</li></ul></div>`,
	)
}

func TestBooleanField(t *testing.T) {
	testBasic(
		t,
		NewForm("BooleanField", Fields(BooleanField("yes", "YES", true), BooleanField("no", "NO", false))),
		`no=no`,
		`<input type="checkbox" name="yes" value="yes" checked >YES<input type="checkbox" name="no" value="no" >NO`,
		`<input type="checkbox" name="yes" value="yes" >YES<input type="checkbox" name="no" value="no" checked >NO`,
	)
}

func TestRadioInput(t *testing.T) {
	testBasic(
		t,
		NewForm("RadioInput", Fields(RadioInput("radio-up", "UP", "up", false), RadioInput("radio-down", "DOWN", "down", false))),
		`radio-up=up`,
		`<input type="radio" name="radio-up" value="up" >UP<input type="radio" name="radio-down" value="down" >DOWN`,
		`<input type="radio" name="radio-up" value="up" checked >UP<input type="radio" name="radio-down" value="down" >DOWN`,
	)
}

func TestCheckBoxInput(t *testing.T) {
	testBasic(
		t,
		NewForm(
			"CheckboxInput",
			Fields(
				CheckboxInput("checkbox-left", "LEFT", "left", false),
				CheckboxInput("checkbox-right", "RIGHT", "right", false),
			),
		),
		`checkbox-left=left&checkbox-right=right`,
		`<input type="checkbox" name="checkbox-left" value="left" >LEFT<input type="checkbox" name="checkbox-right" value="right" >RIGHT`,
		`<input type="checkbox" name="checkbox-left" value="left" checked >LEFT<input type="checkbox" name="checkbox-right" value="right" checked >RIGHT`,
	)
}

func makeselectoptions(so ...string) []*Selection {
	var newso []*Selection
	for _, s := range so {
		newso = append(newso, NewSelection(s, strings.ToUpper(s), false))
	}
	return newso
}

func TestSelectField(t *testing.T) {
	testBasic(
		t,
		NewForm("SelectField", Fields(SelectField("selectfield", makeselectoptions("one", "two", "three"), nil, nil))),
		`selectfield=three`,
		`<select name="selectfield" ><option value="one">ONE</option><option value="two">TWO</option><option value="three">THREE</option></select>`,
		`<select name="selectfield" ><option value="one">ONE</option><option value="two">TWO</option><option value="three" selected>THREE</option></select>`,
	)
	testBasic(
		t,
		NewForm("MultiSelectField", Fields(SelectField("multiselectfield", makeselectoptions("one", "two", "three"), nil, nil, "multiple"))),
		`multiselectfield=one two three`,
		`<select name="multiselectfield" multiple><option value="one">ONE</option><option value="two">TWO</option><option value="three">THREE</option></select>`,
		`<select name="multiselectfield" multiple><option value="one" selected>ONE</option><option value="two" selected>TWO</option><option value="three" selected>THREE</option></select>`,
	)
}

func TestRadioField(t *testing.T) {
	testBasic(
		t,
		NewForm("RadioField", Fields(RadioField("radiofield-group", "Select one:", makeselectoptions("A", "B", "C", "D"), nil, nil))),
		`radiofield-group=A`,
		`<fieldset name="radiofield-group" ><legend>Select one:</legend><ul><li><input type="radio" name="radiofield-group" value="A" >A</li><li><input type="radio" name="radiofield-group" value="B" >B</li><li><input type="radio" name="radiofield-group" value="C" >C</li><li><input type="radio" name="radiofield-group" value="D" >D</li></ul></fieldset>`,
		`<fieldset name="radiofield-group" ><legend>Select one:</legend><ul><li><input type="radio" name="radiofield-group" value="A" checked >A</li><li><input type="radio" name="radiofield-group" value="B" >B</li><li><input type="radio" name="radiofield-group" value="C" >C</li><li><input type="radio" name="radiofield-group" value="D" >D</li></ul></fieldset>`,
	)
}

func TestDateField(t *testing.T) {
	testBasic(
		t,
		NewForm("DateField", Fields(DateField("datefield"))),
		`datefield=26/02/2015`,
		`<input type="date" name="datefield" value="" >`,
		`<input type="date" name="datefield" value="26/02/2015" >`,
	)
	testBasic(
		t,
		NewForm("DateField :: Error", Fields(DateField("datefield"))),
		`datefield=26022015`,
		`<input type="date" name="datefield" value="" >`,
		`<input type="date" name="datefield" value="26022015" ><div class="field-errors"><ul><li>Cannot parse 26022015 in format 02/01/2006</li></ul></div>`,
	)
}

func TestListField(t *testing.T) {
	var ListField1 Field = ListField("listfield", 3, TextField("TEST", nil, nil))
	testBasic(
		t,
		NewForm("ListField", Fields(ListField1)),
		`listfield-0-TEST=IamZERO&listfield-1-TEST=IamONE&listfield7-seven=IshouldnotbeSEVEN`,
		`<fieldset name="listfield" ><ul><li><input type="text" name="listfield-0-TEST" value="" ></li><li><input type="text" name="listfield-1-TEST" value="" ></li><li><input type="text" name="listfield-2-TEST" value="" ></li></ul></fieldset>`,
		`<fieldset name="listfield" ><ul><li><input type="text" name="listfield-0-TEST" value="IamZERO" ></li><li><input type="text" name="listfield-1-TEST" value="IamONE" ></li></ul></fieldset>`,
	)
}

func SimpleForm() Form {
	return NewForm("simple", Fields(TextField("fftext", nil, nil), BooleanField("yes", "Yes", false)))
}

func TestFormsField(t *testing.T) {
	var FormsField1 Field = FormsField("formfield", 1, SimpleForm())

	testBasic(
		t,
		NewForm("FormsField", Fields(FormsField1)),
		`formfield=2&formfield-0-fftext-0=TEXTFIELD0&formfield-0-yes-1=yes&formfield-1-fftext-0=TEXTFIELD1`,
		`<fieldset name="formfield"><input type="hidden" name="formfield" value="1"><ul><li><input type="text" name="formfield-0-fftext-0" value="" ><input type="checkbox" name="formfield-0-yes-1" value="yes" >Yes</li></ul></fieldset>`,
		`<fieldset name="formfield"><input type="hidden" name="formfield" value="2"><ul><li><input type="text" name="formfield-0-fftext-0" value="TEXTFIELD0" ><input type="checkbox" name="formfield-0-yes-1" value="yes" checked >Yes</li><li><input type="text" name="formfield-1-fftext-0" value="TEXTFIELD1" ><input type="checkbox" name="formfield-1-yes-1" value="yes" >Yes</li></ul></fieldset>`,
	)
}

func extracttoken(from *httptest.ResponseRecorder, by string) string {
	b := bytes.Fields(from.Body.Bytes())

	for _, x := range b {
		if bytes.Contains(x, []byte(by)) {
			return strings.Split(string(x), `"`)[1]
		}
	}

	return ""
}

func TestXSRF(t *testing.T) {
	field := XSRF("testXSRF", "SECRET")
	xsrffield, _ := field.(*xsrf)

	f := NewForm("XSRF", Fields(field))

	ts := testserve()
	ts.handlers["GET"] = formHandlerGet(t, f)
	ts.handlers["POST"] = formHandlerPost(t, f)

	w1 := PerformGet(ts)

	expect := `<input type="hidden" name="testXSRF" value="`

	if !strings.Contains(w1.Body.String(), expect) {
		t.Errorf("\n%s XSRF GET Error\ngot\n\n%s\nshould contain:\n\n%s\n\n", w1.Body, expect)
	}

	token := extracttoken(w1, "value=")

	if !validTokenAtTime(token, xsrffield.Secret, xsrffield.Key, time.Now()) {
		t.Errorf("\nInvlaid xsrf token: %s\n", token)
	}

	sendtoken := fmt.Sprintf(`testXSRF=%s`, token)

	w2 := PerformPost(ts, sendtoken)

	posttoken := extracttoken(w2, "value=")

	if posttoken == token {
		t.Errorf("\n %s == %s ; Indicative of an invalid token\n", token, posttoken)
	}

	w3 := PerformPost(ts, `testXSRF=invalidtoken`)

	invalidtoken := extracttoken(w3, "value=")

	if invalidtoken != "invalidtoken" {
		t.Errorf(`\nInvalid token should be "invalidtoken", not %s\n`, invalidtoken)
	}

	var invalidresult string = `<input type="hidden" name="testXSRF" value="invalidtoken" ><div class="field-errors"><ul><li>Invalid XSRF Token</li></ul></div>`

	if !strings.Contains(w3.Body.String(), invalidresult) {
		t.Errorf("\n%s POST Error\ngot %s\nshould contain %s\n\n", w3, invalidresult)
	}
}

func TestSubmitField(t *testing.T) {
	testBasic(
		t,
		NewForm("SubmitField", Fields(SubmitField("test", nil, nil))),
		``,
		`<form action="/" method="POST"><input type="submit" name="test" value="test" ></form>`,
		`<form action="/" method="POST"><input type="submit" name="test" value="test" ></form>`,
	)
}
