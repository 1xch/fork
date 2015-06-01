package fork

import (
	"fmt"
	"reflect"
)

type Value struct {
	Raw  interface{}
	Type reflect.Type
	Kind reflect.Kind
}

func NewValue(i interface{}) *Value {
	t := reflect.TypeOf(i)
	return &Value{
		Raw:  i,
		Type: t,
		Kind: t.Kind(),
	}
}

func (v *Value) String() string {
	if v.Raw != nil {
		switch v.Raw.(type) {
		case *Selection, []*Selection, []Field, []Form:
			return fmt.Sprintf("%+v", v.Raw)
		case bool:
			return fmt.Sprintf("%t", v.Raw)
		case int:
			return fmt.Sprintf("%d", v.Raw)
		case fmt.Stringer:
			return v.String()
		default:
			return fmt.Sprintf("%s", v.Raw)
		}
	}
	return ""
}

func (v *Value) Integer() int {
	if v.Raw != nil {
		switch v.Raw.(type) {
		case *Selection, []*Selection, []Field, Form, []Form:
			return -1
		default:
			return v.Raw.(int)
		}
	}
	return 0
}

func (v *Value) Bool() bool {
	if v.Raw != nil {
		switch v.Raw.(type) {
		case *Selection:
			return v.Raw.(*Selection).Set
		case []*Selection, []Field, Form, []Form:
			return false
		default:
			return v.Raw.(bool)
		}
	}
	return false
}
