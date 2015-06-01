package fork

import (
	"fmt"
	"reflect"
)

type Validater interface {
	Error(Field) bool
	Errors(Field) []string
	Valid(Field) bool
	Validate(Field) error
}

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

func (v *validater) AddValidater(fn interface{}) {
	v.validaters = append(v.validaters, valueValidater(fn))
}

func reflectValidaters(fns ...interface{}) []reflect.Value {
	var ret []reflect.Value
	for _, fn := range fns {
		ret = append(ret, valueValidater(fn))
	}
	return ret
}

func valueValidater(fn interface{}) reflect.Value {
	v := reflect.ValueOf(fn)
	if !isValidater(v.Type()) {
		panic(fmt.Sprintf("Cannot use %q as a validater function: function must return 1 value and the value must be an error", fn))
	}
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("%+v is not a function", fn))
	}
	return v
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
