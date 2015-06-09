package fork

import "reflect"

type Checker interface {
	Valid(Form) bool
	Check(Form) (bool, error)
	Checks(...interface{}) []reflect.Value
}

type checker struct {
	checked bool
	checks  []reflect.Value
}

func NewChecker() Checker {
	return &checker{}
}

func (c *checker) Valid(f Form) bool {
	for _, fd := range f.Fields() {
		if !fd.Valid(fd) {
			return false
		}
	}
	return true
}

func (c *checker) Check(f Form) (bool, error) {
	var err error
	var valid bool
	for _, fn := range c.checks {
		valid, err = Check(fn, f)
		if err != nil {
			c.checked = true
			return valid, err
		}
	}
	c.checked = true
	return valid, err
}

func (c *checker) Checks(fns ...interface{}) []reflect.Value {
	for _, fn := range fns {
		c.AddCheck(fn)
	}
	return c.checks
}

func (c *checker) AddCheck(fn interface{}) {
	c.checks = append(c.checks, valueFn(fn, isCheck, `must return 1 boolean and 1 error`))
}

var boolType = reflect.TypeOf((*bool)(nil)).Elem()

func isCheck(typ reflect.Type) bool {
	switch {
	case typ.NumOut() == 2 && typ.Out(0) == boolType && typ.Out(1) == errorType:
		return true
	}
	return false
}

var BadCheck = ForkError(`check function did not return a boolean value with its error.`)

func Check(fn reflect.Value, args ...interface{}) (bool, error) {
	checked, err := call(fn, args...)
	if checked, ok := checked.(bool); ok {
		return checked, err
	}
	return false, BadCheck
}
