package fork

import (
	"fmt"
	"net/http"
	"reflect"
)

type Filterer interface {
	Filter(string, *http.Request) *Value
}

func NewFilterer(f ...interface{}) Filterer {
	return &filterer{filters: reflectFilters(f...)}
}

type filterer struct {
	filters []reflect.Value
}

func (fr *filterer) Filter(k string, r *http.Request) *Value {
	var v interface{}
	v = r.FormValue(k)
	for _, fn := range fr.filters {
		v = Filter(fn, v)
	}
	if v != nil {
		return NewValue(v)
	}
	return NewValue(nil)
}

func (fr *filterer) AddFilter(fn interface{}) {
	fr.filters = append(fr.filters, valueFilter(fn))
}

func reflectFilters(fns ...interface{}) []reflect.Value {
	var ret []reflect.Value
	for _, fn := range fns {
		ret = append(ret, valueFilter(fn))
	}
	return ret
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func valueFilter(fn interface{}) reflect.Value {
	v := reflect.ValueOf(fn)
	if !isFilter(v.Type()) {
		panic(fmt.Sprintf("Cannot use function %q with %d results\nreturn must be 1 value, or 1 value and 1 error value", fn, v.Type().NumOut()))
	}
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("%+v is not a function", fn))
	}
	return v
}

func isFilter(typ reflect.Type) bool {
	switch {
	case typ.NumOut() == 1:
		return true
	case typ.NumOut() == 2 && typ.Out(1) == errorType:
		return true
	}
	return false
}

func Filter(fn reflect.Value, args ...interface{}) interface{} {
	filtered, err := call(fn, args...)
	if err != nil {
		return err
	}
	return filtered
}
