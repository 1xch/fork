package fork

import "testing"

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
