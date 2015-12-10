package fork

import (
	"strings"
	"testing"
)

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
