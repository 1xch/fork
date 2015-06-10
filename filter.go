package fork

import (
	"net/http"
	"reflect"
)

type Filterer interface {
	Filter(string, *http.Request) *Value
	Filters(...interface{}) []reflect.Value
}

func NewFilterer(f ...interface{}) Filterer {
	return &filterer{filters: reflectFilters(f...)}
}

var nilFilterer = &filterer{}

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

func (fr *filterer) Filters(fns ...interface{}) []reflect.Value {
	fr.filters = append(fr.filters, reflectFilters(fns...)...)
	return fr.filters
}

func reflectFilters(fns ...interface{}) []reflect.Value {
	var ret []reflect.Value
	for _, fn := range fns {
		ret = append(ret, valueFn(fn, isFilter, `1 value, or 1 value and 1 error value`))
	}
	return ret
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
