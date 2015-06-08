package fork

import "reflect"

type Validater interface {
	Error(Field) bool
	Errors(Field) []string
	Valid(Field) bool
	Validate(Field) error
	Validaters(...interface{}) []reflect.Value
}

var nilValidater = &validater{}

func NewValidater(v ...interface{}) Validater {
	return &validater{validaters: reflectValidaters(v...)}
}

type validater struct {
	validaters []reflect.Value
}

func (v *validater) Error(f Field) bool {
	return !v.Valid(f)
}

func (v *validater) Errors(f Field) []string {
	var ret []string
	for _, vdr := range v.validaters {
		err := Validate(vdr, f)
		if err != nil {
			ret = append(ret, err.Error())
		}
	}
	return ret
}

func (v *validater) Valid(f Field) bool {
	err := v.Validate(f)
	if err != nil {
		return false
	}
	return true
}

func (v *validater) Validate(f Field) error {
	for _, vdr := range v.validaters {
		err := Validate(vdr, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *validater) Validaters(fns ...interface{}) []reflect.Value {
	v.validaters = append(v.validaters, reflectValidaters(fns...)...)
	return v.validaters
}

func reflectValidaters(fns ...interface{}) []reflect.Value {
	var ret []reflect.Value
	for _, fn := range fns {
		ret = append(ret, valueFn(fn, `1 value and the value must be an error`))
	}
	return ret
}

func isValidater(typ reflect.Type) bool {
	switch {
	case typ.NumOut() == 1 && typ.Out(0) == errorType:
		return true
	}
	return false
}

func Validate(fn reflect.Value, args ...interface{}) error {
	validated, err := call(fn, args...)
	if err != nil {
		return err
	}
	if validated != nil {
		return validated.(error)
	}
	return nil
}
