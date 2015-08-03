package fork

import (
	"fmt"
	"reflect"
)

type forkError struct {
	err  string
	vals []interface{}
}

func (f *forkError) Error() string {
	return fmt.Sprintf("[fork] %s", fmt.Sprintf(f.err, f.vals...))
}

func (f *forkError) Out(vals ...interface{}) *forkError {
	f.vals = vals
	return f
}

func ForkError(err string) *forkError {
	return &forkError{err: err}
}

var NotAFunction = ForkError(`#+v is not a function`).Out

var InvalidFunction = ForkError(`cannot use function %q with %d results, return must be %s`).Out

var WrongNumberArgs = ForkError(`wrong number of args: got %d want at least %d`).Out

var UnassignableArg = ForkError(`arg %d has type %s; should be %s`).Out

func valueFn(fn interface{}, is func(reflect.Type) bool, out string) reflect.Value {
	v := reflect.ValueOf(fn)
	if !is(v.Type()) {
		panic(InvalidFunction(fn, v.Type().NumOut(), out).Error())
	}
	if v.Kind() != reflect.Func {
		panic(NotAFunction(fn).Error())
	}
	return v
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

var boolType = reflect.TypeOf((*bool)(nil)).Elem()

func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}

func call(fn reflect.Value, args ...interface{}) (interface{}, error) {
	typ := fn.Type()
	numIn := typ.NumIn()
	var dddType reflect.Type
	if typ.IsVariadic() {
		if len(args) < numIn-1 {
			return nil, WrongNumberArgs(len(args), numIn-1)
		}
		dddType = typ.In(numIn - 1).Elem()
	} else {
		if len(args) != numIn {
			return nil, WrongNumberArgs(len(args), numIn)
		}
	}
	argv := make([]reflect.Value, len(args))
	for i, arg := range args {
		value := reflect.ValueOf(arg)
		var argType reflect.Type
		if !typ.IsVariadic() || i < numIn-1 {
			argType = typ.In(i)
		} else {
			argType = dddType
		}
		if !value.IsValid() && canBeNil(argType) {
			value = reflect.Zero(argType)
		}
		if !value.Type().AssignableTo(argType) {
			return nil, UnassignableArg(i, value.Type(), argType)
		}
		argv[i] = value
	}
	result := fn.Call(argv)
	if len(result) == 2 && !result[1].IsNil() {
		return result[0].Interface(), result[1].Interface().(error)
	}
	return result[0].Interface(), nil
}
