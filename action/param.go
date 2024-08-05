package action

import (
	"fmt"
	"reflect"
)

type ActionParam interface {
	GetType() string
	GetValue() interface{}
}

type ReflectActionParam struct {
	value interface{}
	typ   reflect.Type
}

func (ra *ReflectActionParam) GetType() string {
	if ra.typ == nil {
		ra.typ = reflect.TypeOf(ra.value)
	}

	return ra.typ.String()
}

func (ra *ReflectActionParam) GetValue() interface{} {
	return ra.value
}

func NewParamReflect(value interface{}) *ReflectActionParam {
	return &ReflectActionParam{
		value: value,
	}
}

type SimpleActionParam struct {
	Type  string
	Value interface{}
}

func (sp *SimpleActionParam) GetType() string {
	return sp.Type
}

func (sp *SimpleActionParam) GetValue() interface{} {
	return sp.Value
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

type MargedParam []ActionParam

func (mp MargedParam) GetValue() interface{} {
	return mp[0].GetValue()
}

func (mp MargedParam) GetType() string {
	return mp[0].GetType()
}
