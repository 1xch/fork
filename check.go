package fork

import "reflect"

// Checker is a form level representation of functions added to the form that
// may range over the entire form to assess the form as a whole, useful for comparing
// or assessing multiple fields through one or any number of functions.
//
// A check function is any function that accepts a form and returns a boolean
// an error.
type Checker interface {
	Checkable(Form) bool
	Check(Form) error
	MustCheck(Form) error
	Error(Form) bool
	Errors(Form) []string
	Checks(...interface{}) []reflect.Value
}

type checker struct {
	checks []reflect.Value
}

// NewChecker returns a new default Checker.
func NewChecker(c ...interface{}) Checker {
	return &checker{checks: reflectChecks(c...)}
}

// Checkable returns whether the form is in a state to run checks or not,
// optionally ensuring all fields are valid before marking the form as checkable.
func (c *checker) Checkable(f Form) bool {
	for _, fd := range f.Fields() {
		if !fd.Valid(fd) {
			return false
		}
	}
	return true
}

// Check runs all form checks, regardless of form checked state or field validity.
func (c *checker) Check(f Form) error {
	for _, fn := range c.checks {
		err := Check(fn, f)
		if err != nil {
			return err
		}
	}
	return nil
}

var NotCheckableError = Frror("Form is not checkable, due to invalid fields.")

// Mustcheck ensures that the form is ok(all fields are valid and without error)
// to check before running all form checks.
func (c *checker) MustCheck(f Form) error {
	checkable := c.Checkable(f)
	if checkable {
		err := c.Check(f)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
	return NotCheckableError
}

// Error provides a boolean indicating if the form contains any type of check
// or field validation errors.
func (c *checker) Error(f Form) bool {
	if len(c.Errors(f)) > 0 {
		return true
	}
	return false
}

// If the form contains errors of Checks, or field errors, the length of
// the returned string slice will be greater than zero.
func (c *checker) Errors(f Form) []string {
	var ret []string
	checkable := c.Checkable(f)
	if !checkable {
		ret = append(ret, NotCheckableError.Error())
	}
	if checkable {
		for _, fn := range c.checks {
			err := Check(fn, f)
			if err != nil {
				ret = append(ret, err.Error())
			}
		}
	}
	for _, fd := range f.Fields() {
		ret = append(ret, fd.Errors(fd)...)
	}
	return ret
}

func reflectChecks(fns ...interface{}) []reflect.Value {
	var ret []reflect.Value
	for _, fn := range fns {
		ret = append(ret, valueFn(fn, isCheck, `must 1 error`))
	}
	return ret
}

// Checks takes any number of Check functions add them to and returns a slice
// of the reflect.Value functions associated with this Checker. Fork will panic
// if the functions passed in are not acceptable Check functions(i.e. providing
// a signature func(Form) error)
func (c *checker) Checks(fns ...interface{}) []reflect.Value {
	c.checks = reflectChecks(fns...)
	return c.checks
}

func isCheck(typ reflect.Type) bool {
	if typ.NumOut() == 1 && typ.Out(0) == errorType {
		return true
	}
	return false
}

// Check is a function that takes any reflect.value function, any varadic number
// of arguments, and returns a boolean and an error.
func Check(fn reflect.Value, args ...interface{}) error {
	ckd, err := call(fn, args...)
	if err != nil {
		return err
	}
	var ret error
	switch ckd.(type) {
	case error:
		ret = ckd.(error)
	case nil:
		ret = nil
	}
	return ret
}
