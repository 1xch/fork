package fork

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func formstring(note string, f Form) *bytes.Buffer {
	out := bytes.NewBuffer([]byte(fmt.Sprintf("begin: %s", note)))
	out.WriteString(`<form action="/" method="POST">`)
	out.WriteString(f.String())
	out.WriteString("</form>")
	out.WriteString(fmt.Sprintf("end: %s", note))
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
		out := formstring(`form via GET`, t)
		w.Write(out.Bytes())
	}
}

func postformhandlerfor(f Form) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		t := f.New()
		t.Process(r)
		out := formstring(`form via POST`, t)
		w.Write(out.Bytes())
	}
}

func testbasic(t *testing.T, name string, f Form, postprovides string, GETexpects string, POSTexpects string) {
	w1, w2 := PerformForForm(f, postprovides)

	if w1.Code != 200 || w2.Code != 200 {
		t.Errorf("Response incorrect; received Get %d Post d%, expected 200", w1.Code, w2.Code)
	}

	if !strings.Contains(w1.Body.String(), GETexpects) {
		t.Errorf("%s GET Error\ngot\n\n%s\nshould contain:\n\n%s\n\n", name, w1.Body, GETexpects)
	}

	if !strings.Contains(w2.Body.String(), POSTexpects) {
		t.Errorf("%s POST Error\ngot %s\nshould contain %s\n\n", name, w2.Body, POSTexpects)
	}
}

func TestTextField(t *testing.T) {
	testbasic(
		t,
		"TextField",
		NewForm(TextField("text")),
		`text=TEXT`,
		`<input type="text" name="text" value="" >`,
		`<input type="text" name="text" value="TEXT" >`,
	)
}

func TestTextAreaField(t *testing.T) {
	testbasic(
		t,
		"TextAreaField",
		NewForm(TextAreaField("textarea", "rows=10", "cols=10")),
		`textarea=TEXTAREA`,
		`<textarea name="textarea" rows=10 cols=10></textarea>`,
		`<textarea name="textarea" rows=10 cols=10>TEXTAREA</textarea>`,
	)
}

func TestHiddenField(t *testing.T) {
	testbasic(
		t,
		"HiddenField",
		NewForm(HiddenField("hidden")),
		`hidden=HIDDEN`,
		`<input type="hidden" name="hidden" value="" >`,
		`<input type="hidden" name="hidden" value="HIDDEN" >`,
	)
}

func TestPasswordField(t *testing.T) {
	testbasic(
		t,
		"PasswordField",
		NewForm(PassWordField("password", "size=10", "maxlength=30")),
		`password=PASSWORD`,
		`<input type="password" name="password" value="" size=10 maxlength=30>`,
		`<input type="password" name="password" value="PASSWORD" size=10 maxlength=30>`,
	)
}

func TestBooleanField(t *testing.T) {
	testbasic(
		t,
		"BooleanField",
		NewForm(BooleanField("yes", "YES", true), BooleanField("no", "NO", false)),
		`no=no`,
		`<input type="checkbox" name="yes" value="yes" checked >YES<input type="checkbox" name="no" value="no" >NO`,
		`<input type="checkbox" name="yes" value="yes" >YES<input type="checkbox" name="no" value="no" checked >NO`,
	)
}

func TestRadioInput(t *testing.T) {
	testbasic(
		t,
		"RadioInput",
		NewForm(RadioInput("radio-up", "UP", "up", false), RadioInput("radio-down", "DOWN", "down", false)),
		`radio-up=up`,
		`<input type="radio" name="radio-up" value="up" >UP<input type="radio" name="radio-down" value="down" >DOWN`,
		`<input type="radio" name="radio-up" value="up" checked >UP<input type="radio" name="radio-down" value="down" >DOWN`,
	)
}

func TestCheckBoxInput(t *testing.T) {
	testbasic(
		t,
		"CheckboxInput",
		NewForm(CheckboxInput("checkbox-left", "LEFT", "left", false), CheckboxInput("checkbox-right", "RIGHT", "right", false)),
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
	testbasic(
		t,
		"SelectField",
		NewForm(SelectField("selectfield", makeselectoptions("one", "two", "three")...)),
		`selectfield=three`,
		`<select name="selectfield" ><option value="one">ONE</option><option value="two">TWO</option><option value="three">THREE</option></select>`,
		`<select name="selectfield" ><option value="one">ONE</option><option value="two">TWO</option><option value="three" selected>THREE</option></select>`,
	)
}

func TestMultiSelectField(t *testing.T) {
	testbasic(
		t,
		"MultiSelectField",
		NewForm(MultiSelectField("multiselectfield", makeselectoptions("one", "two", "three")...)),
		`multiselectfield=one two three`,
		`<select name="multiselectfield" multiple><option value="one">ONE</option><option value="two">TWO</option><option value="three">THREE</option></select>`,
		`<select name="multiselectfield" multiple><option value="one" selected>ONE</option><option value="two" selected>TWO</option><option value="three" selected>THREE</option></select>`,
	)
}

func TestRadioField(t *testing.T) {
	testbasic(
		t,
		"RadioField",
		NewForm(RadioField("radiofield-group", "Select one:", makeselectoptions("A", "B", "C", "D")...)),
		`radiofield-group=A`,
		`<fieldset name="radiofield-group"><legend>Select one:</legend><ul><li><input type="radio" name="radiofield-group" value="A" >A</li><li><input type="radio" name="radiofield-group" value="B" >B</li><li><input type="radio" name="radiofield-group" value="C" >C</li><li><input type="radio" name="radiofield-group" value="D" >D</li></ul></fieldset>`,
		`<fieldset name="radiofield-group"><legend>Select one:</legend><ul><li><input type="radio" name="radiofield-group" value="A" checked >A</li><li><input type="radio" name="radiofield-group" value="B" >B</li><li><input type="radio" name="radiofield-group" value="C" >C</li><li><input type="radio" name="radiofield-group" value="D" >D</li></ul></fieldset>`,
	)
}

func TestListField(t *testing.T) {
	var ListField1 Field = ListField("listfield", 3, TextField("TEST"))

	testbasic(
		t,
		"ListField",
		NewForm(ListField1),
		`listfield-0-TEST=IamZERO&listfield-1-TEST=IamONE&listfield7-seven=IshouldnotbeSEVEN`,
		`<fieldset name="listfield" ><ul><li><input type="text" name="listfield-0-TEST" value="" ></li><li><input type="text" name="listfield-1-TEST" value="" ></li><li><input type="text" name="listfield-2-TEST" value="" ></li></ul></fieldset>`,
		`<fieldset name="listfield" ><ul><li><input type="text" name="listfield-0-TEST" value="IamZERO" ></li><li><input type="text" name="listfield-1-TEST" value="IamONE" ></li></ul></fieldset>`,
	)
}

func SimpleForm() Form {
	return NewForm(TextField("fftext"), BooleanField("yes", "Yes", false))
}

func TestFormsField(t *testing.T) {
	var FormsField1 Field = FormsField("formfield", 1, SimpleForm())

	testbasic(
		t,
		"FormsField",
		NewForm(FormsField1),
		`formfield=2&formfield-0-fftext-0=TEXTFIELD0&formfield-0-yes-1=yes&formfield-1-fftext-0=TEXTFIELD1`,
		`<fieldset name="formfield"><input type="hidden" name="formfield" value="1"><ul><li><input type="text" name="formfield-0-fftext-0" value="" ><input type="checkbox" name="formfield-0-yes-1" value="yes" >Yes</li></ul></fieldset>`,
		`<fieldset name="formfield"><input type="hidden" name="formfield" value="2"><ul><li><input type="text" name="formfield-0-fftext-0" value="TEXTFIELD0" ><input type="checkbox" name="formfield-0-yes-1" value="yes" checked >Yes</li><li><input type="text" name="formfield-1-fftext-0" value="TEXTFIELD1" ><input type="checkbox" name="formfield-1-yes-1" value="yes" >Yes</li></ul></fieldset>`,
	)
}
