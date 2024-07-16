package action

import "fmt"

type ActionParam interface {
	GetType() string
	GetValue() interface{}
}

func As[T any](parm ActionParam) (T, error) {
	v := parm.GetValue()

	if v != nil {
		if val, ok := v.(T); ok {
			return val, nil
		}
	}
	var t T

	return t, fmt.Errorf("type assertion failed: %v is not %T", v, t)
}
