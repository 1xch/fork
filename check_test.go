package fork

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type checkExpectation struct {
	name, fill string
	expects    error
	fn         func(Form) error
}

func basicCheck(Form) error {
	return nil
}

var basicCheckExpectation = checkExpectation{"basicCheck", "basic", nil, basicCheck}

var erredError = errors.New("There was an error")

func erredCheck(Form) error {
	return erredError
}

var erredCheckExpectation = checkExpectation{"erredCheck", "erred", erredError, erredCheck}

var errFieldLenghtIsNOTOne = errors.New("fields length is not equal to 1")

var errFieldLengthIsOne = errors.New("field length is one")

func arbitraryCheck(f Form) error {
	if len(f.Fields()) != 1 {
		return errFieldLenghtIsNOTOne
	}
	return errFieldLengthIsOne
}

var arbitraryCheckExpectation = checkExpectation{"arbCheck", "arb", errFieldLengthIsOne, arbitraryCheck}

var notCheckableExpectation = checkExpectation{"notCheckable", "noCheck", NotCheckableError, arbitraryCheck}

func checkError(t *testing.T, f Form, ce checkExpectation) {
	if ce.expects != nil {
		if !f.Error(f) {
			t.Error("form.Error() expected to be true, but was not")
		}
		errs := f.Errors(f)
		expe := ce.expects.Error()
		var contains bool
		for _, v := range errs {
			if strings.Contains(v, expe) {
				contains = true
			}
		}
		if !contains {
			t.Errorf("Expected %s in %s", expe, errs)
		}
	} else {
		if f.Error(f) {
			t.Error("Form reported errors, but there is no expectation of error.")
		}
	}
}

type testCheckFunc func(t *testing.T, f Form, ce checkExpectation)

func testCheck(t *testing.T, f Form, ce checkExpectation) {
	if err := f.Check(f); ce.expects != err {
		t.Errorf("Check error, expected %s, got %s", ce.expects, err)
	}
}

func testMustCheck(t *testing.T, f Form, ce checkExpectation) {
	if err := f.MustCheck(f); ce.expects != err {
		t.Errorf("MustCheck error, expected %s, got %s", ce.expects, err)
	}
}

func formCheckHandlerPost(t *testing.T, f Form, ce checkExpectation, fn testCheckFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		nf := f.New()
		nf.Process(r)
		fn(t, f, ce)
		checkError(t, f, ce)
		w.WriteHeader(200)
		out := wrapForm(nf)
		w.Write(out.Bytes())
	}
}

func validForm(fname string, check func(Form) error) Form {
	return NewForm("Test Check", Fields(HiddenField(fname, nil, nil)), Checks(check))
}

func invalidateField(f *textField) error {
	if f.Validateable() {
		return errors.New("invalid field")
	}
	return nil
}

func invalidForm(fname string, check func(Form) error) Form {
	nf := NewForm("Test Check invalid form", Fields(HiddenField(fname, []interface{}{invalidateField}, nil)))
	nf.Checks(check)
	return nf
}

func checkTest(t *testing.T, valid bool, ce checkExpectation, must bool) {
	var form Form
	if valid {
		form = validForm(ce.name, ce.fn)
	} else {
		form = invalidForm(ce.name, ce.fn)
	}
	ts := testServe()
	ts.handlers["GET"] = formHandlerGet(t, form)
	if must {
		ts.handlers["POST"] = formCheckHandlerPost(t, form, ce, testMustCheck)
	} else {
		ts.handlers["POST"] = formCheckHandlerPost(t, form, ce, testCheck)
	}
	PerformGet(ts)
	PerformPost(ts, fmt.Sprintf(`%s=%s`, ce.name, ce.fill))
}

func TestChecks(t *testing.T) {
	checkTest(t, true, basicCheckExpectation, true)
	checkTest(t, true, erredCheckExpectation, true)
	checkTest(t, true, arbitraryCheckExpectation, true)
	checkTest(t, true, basicCheckExpectation, false)
	checkTest(t, true, erredCheckExpectation, false)
	checkTest(t, true, arbitraryCheckExpectation, false)
	checkTest(t, false, notCheckableExpectation, true)
}
