package action

import (
	"fmt"
	"reflect"
)

type Action interface {
	Handle(params ...ActionParam) ([]ActionParam, error)
	InputTypes() []string
	OutputTypes() []string
}

type FuncAction struct {
	fn           interface{}
	fnType       reflect.Type
	fnValue      reflect.Value
	inputTypes   []string
	inputParsed  bool
	outputTypes  []string
	outputParsed bool
}

func (f *FuncAction) Handle(params ...ActionParam) ([]ActionParam, error) {
	if f.fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("action error: not a function")
	}

	lenParam := len(params)
	in := make([]reflect.Value, lenParam)
	numIn := f.fnType.NumIn()
	if numIn != lenParam {
		return nil, fmt.Errorf("invalid number of params: expected %d parameters, %d given", numIn, lenParam)
	}
	inTypes := f.InputTypes()

	for i, param := range params {
		if inTypes[i] != param.GetType() {
			return nil, fmt.Errorf("invalid type of param[%d]: expected %s, %s given", i, inTypes[i], param.GetType())
		}

		in[i] = reflect.ValueOf(param.GetValue())
	}

	values := f.fnValue.Call(in)

	out := make([]ActionParam, len(values))

	for i, value := range values {
		out[i] = &SimpleActionParam{
			Type:  value.Type().String(),
			Value: value.Interface(),
		}
	}

	return out, nil
}

func (f *FuncAction) InputTypes() (types []string) {
	if f.inputParsed {
		return f.inputTypes
	}

	if f.fnType.Kind() != reflect.Func {
		return
	}

	for i := 0; i < f.fnType.NumIn(); i++ {
		pt := f.fnType.In(i)
		types = append(types, pt.String())
	}

	f.inputTypes = types
	f.inputParsed = true

	return
}

func (f *FuncAction) OutputTypes() (types []string) {
	if f.outputParsed {
		return f.outputTypes
	}
	if f.fnType.Kind() != reflect.Func {
		return
	}

	for i := 0; i < f.fnType.NumOut(); i++ {
		pt := f.fnType.Out(i)
		types = append(types, pt.String())
	}

	f.outputTypes = types
	f.outputParsed = true
	return
}

func Func(fn interface{}) *FuncAction {

	val := reflect.ValueOf(fn)

	return &FuncAction{
		fn:      fn,
		fnType:  val.Type(),
		fnValue: val,
	}
}
