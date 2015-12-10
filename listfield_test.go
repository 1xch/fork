package fork

import "testing"

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
