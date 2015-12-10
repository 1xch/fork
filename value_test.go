package fork

import (
	"fmt"
	"net/http"
	"testing"
)

func testValueString(t *testing.T, expected, provided string) {
	if provided != expected {
		t.Errorf("Expected %s, but received %s from value.String()", expected, provided)
	}
}

func testValueInt(t *testing.T, expected, provided int) {
	if provided != expected {
		t.Errorf("Expected %d, but received %d from value.Integer()", expected, provided)
	}
}

func testValueBool(t *testing.T, expected, provided bool) {
	if provided != expected {
		t.Errorf("Expected %t, but received %t from value.Boolean()", expected, provided)
	}
}

var (
	nilValue    = &Value{nil}
	testStrings = []*Value{
		&Value{"1"},
		&Value{true},
		&Value{1},
	}
	testIntegers = []*Value{
		&Value{"1"},
		&Value{true},
		&Value{1},
	}
	testBools = []*Value{
		&Value{"true"},
		&Value{true},
		&Value{1},
	}
	valFalse     = &Value{false}
	valArbitrary = &Value{1000}
	valMap       = &Value{new(map[string]int)}
)

func TestRaw(t *testing.T) {
	testValueString(t, "", nilValue.String())
	testValueInt(t, 0, nilValue.Integer())
	testValueBool(t, false, nilValue.Boolean())
	testValueString(t, "1", testStrings[0].String())
	testValueString(t, "true", testStrings[1].String())
	testValueString(t, "1", testStrings[2].String())
	for _, v := range testIntegers {
		testValueInt(t, 1, v.Integer())
	}
	testValueInt(t, 0, valFalse.Integer())
	testValueInt(t, -1, valMap.Integer())
	for _, v := range testBools {
		testValueBool(t, true, v.Boolean())
	}
	testValueBool(t, false, valArbitrary.Boolean())
	testValueBool(t, false, valMap.Boolean())
}

func TestValues(t *testing.T) {
	f := NewForm(
		"TestValues",
		Fields(
			TextField("ValueOne", nil, nil),
			TextField("ValueTwo", nil, nil),
			BooleanField("ValueBool", "VALUEBOOL", true),
			TextField("ValueInt", nil, nil),
		),
	)
	v1, v2, v3, v4 := "VALUE1", "VALUE2", "no", "1"
	post := fmt.Sprintf(`ValueOne=%s&ValueTwo=%s&ValueBool=%s&ValueInt=%s`, v1, v2, v3, v4)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		ff := f.New()
		ff.Process(r)
		v := ff.Values()
		testValueString(t, "VALUE1", v["ValueOne"].String())
		testValueString(t, "VALUE2", v["ValueTwo"].String())
		testValueBool(t, false, v["ValueBool"].Boolean())
		testValueInt(t, 1, v["ValueInt"].Integer())
		out := wrapForm(ff)
		w.Write(out.Bytes())
	}
	ts := testServe()
	ts.handlers["POST"] = handler
	PerformPost(ts, post)
}
