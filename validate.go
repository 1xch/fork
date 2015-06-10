package fork

import "reflect"

type Validater interface {
	Validateable() bool
	SetValidateable(bool) bool
	Validate(Field) error
	Valid(Field) bool
	Error(Field) bool
	Errors(Field) []string
	Validaters(...interface{}) []reflect.Value
}

var nilValidater = &validater{}

type validater struct {
	validateable bool
	validaters   []reflect.Value
}

func NewValidater(v ...interface{}) Validater {
	return &validater{validaters: reflectValidaters(v...)}
}

func (v *validater) Validateable() bool {
	return v.validateable
}

func (v *validater) SetValidateable(b bool) bool {
	v.validateable = b
	return v.Validateable()
}

func (v *validater) Validate(f Field) error {
	var err error
	if v.validateable {
		for _, vdr := range v.validaters {
			err = Validate(vdr, f)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (v *validater) Valid(f Field) bool {
	var valid = true
	if v.validateable {
		err := v.Validate(f)
		if err != nil {
			valid = false
		}
	}
	return valid
}

func (v *validater) Error(f Field) bool {
	if v.validateable {
		return !v.Valid(f)
	}
	return false
}

func (v *validater) Errors(f Field) []string {
	var ret []string
	if v.validateable {
		for _, vdr := range v.validaters {
			err := Validate(vdr, f)
			if err != nil {
				ret = append(ret, err.Error())
			}
		}
	}
	return ret
}

func (v *validater) Validaters(fns ...interface{}) []reflect.Value {
	v.validaters = append(v.validaters, reflectValidaters(fns...)...)
	return v.validaters
}

func reflectValidaters(fns ...interface{}) []reflect.Value {
	var ret []reflect.Value
	for _, fn := range fns {
		ret = append(ret, valueFn(fn, isValidater, `1 value and the value must be an error`))
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
