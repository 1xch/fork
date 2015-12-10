package fork

import "testing"

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
