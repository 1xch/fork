package fork

import "reflect"

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
		return v.Raw.(string)
	}
	return ""
}

func (v *Value) Integer() int {
	if v.Raw != nil {
		return v.Raw.(int)
	}
	return 0
}
