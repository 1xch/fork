package fork

import (
	"fmt"
	"strconv"
)

// Value is a raw interface{} useful for managing form and field data.
type Value struct {
	Raw interface{}
}

func NewValue(i interface{}) *Value {
	return &Value{
		Raw: i,
	}
}

// String returns a string of the Raw interface{} of the Value. It will format
// ints and bools to strings, but otherwise will attempt to format to string.
func (v *Value) String() string {
	if v.Raw != nil {
		switch v.Raw.(type) {
		case bool:
			return fmt.Sprintf("%t", v.Raw)
		case int:
			return fmt.Sprintf("%d", v.Raw)
		default:
			return fmt.Sprintf("%s", v.Raw)
		}
	}
	return ""
}

// Integer attempts to return the Raw interface{} value as an int. If the Raw
// interface{} is not string, bool, or int or otherwise impossible to cast to int,
// the value returned is -1.
func (v *Value) Integer() int {
	if v.Raw != nil {
		switch v.Raw.(type) {
		case string:
			str := v.Raw.(string)
			if i, err := strconv.ParseInt(str, 10, 64); err == nil {
				return int(i)
			}
		case bool:
			b := v.Raw.(bool)
			if b {
				return 1
			}
			return 0
		case int:
			return v.Raw.(int)
		default:
			return -1
		}
	}
	return 0
}

// Boolean attempts to return a bool value from the Value.Raw interface{}.
// If this is not possible, bool value false is returned.
func (v *Value) Boolean() bool {
	if v.Raw != nil {
		switch v.Raw.(type) {
		case bool:
			return v.Raw.(bool)
		case string:
			b := v.Raw.(string)
			if bl, err := strconv.ParseBool(b); err == nil {
				return bl
			}
		case int:
			i := v.Raw.(int)
			if i == 1 {
				return true
			} else {
				return false
			}
		case *Selection:
			return v.Raw.(*Selection).Set
		default:
			return false
		}
	}
	return false
}
