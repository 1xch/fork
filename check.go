package fork

import "reflect"

type Checker interface {
	Checkable() bool
	SetCheckable(bool) bool
	Check(Form) (bool, error)
	Ok(Form) bool
	Error(f Form) bool
	Errors(f Form) []string
	Checks(...interface{}) []reflect.Value
}

type checker struct {
	checkable bool
	checks    []reflect.Value
}

func NewChecker(c ...interface{}) Checker {
	return &checker{checks: reflectChecks(c...)}
}

func (c *checker) Checkable() bool {
	return c.checkable
}

func (c *checker) SetCheckable(b bool) bool {
	c.checkable = b
	return c.Checkable()
}

func (c *checker) Check(f Form) (bool, error) {
	ok := c.Ok(f)
	var err error
	if c.checkable {
		for _, fn := range c.checks {
			ok, err = Check(fn, f)
			if err != nil {
				return ok, err
			}
		}
	}
	return ok, err
}

func (c *checker) Ok(f Form) bool {
	for _, fd := range f.Fields() {
		if !fd.Valid(fd) {
			return false
		}
	}
	c.checkable = true
	return true
}

func (c *checker) Error(f Form) bool {
	//if v. {
	//	return !v.Valid(f)
	//}
	return false
}

func (c *checker) Errors(f Form) []string {
	var ret []string
	//if v.validateable {
	//	for _, vdr := range v.validaters {
	//		err := Validate(vdr, f)
	//		if err != nil {
	//			ret = append(ret, err.Error())
	//		}
	//	}
	//}
	return ret
}

func reflectChecks(fns ...interface{}) []reflect.Value {
	var ret []reflect.Value
	for _, fn := range fns {
		ret = append(ret, valueFn(fn, isCheck, `must return 1 boolean and 1 error`))
	}
	return ret
}

func (c *checker) Checks(fns ...interface{}) []reflect.Value {
	c.checks = reflectChecks(fns...)
	return c.checks
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
