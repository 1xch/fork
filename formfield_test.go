package fork

import "testing"

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
