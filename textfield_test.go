package fork

import "testing"

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
