package fork

import "testing"

func TestSubmitField(t *testing.T) {
	testBasic(
		t,
		NewForm("SubmitField", Fields(SubmitField("submit", nil, nil))),
		``,
		`<form action="/" method="POST"><input type="submit" name="submit" value="submit" ><input type="text" name="test" value="" ></form>`,
		`<form action="/" method="POST"><input type="submit" name="submit" value="submit" ><input type="text" name="test" value="FILTERED" ></form>`,
	)
}
