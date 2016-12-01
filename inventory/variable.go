package inventory

import (
	"fmt"
	"encoding/json"
)

const (
	VariableCastString = "string"
	VariableCastInt = "int"
	VariableCastFloat = "float"
	VariableCastBool = "bool"
	VariableCastJson = "json"
)


// Variable represents zero-level inventory variable
type Variable struct {

	// Variable name
	Name string

	// Variable raw value
	Value string

	// Variable output cast
	CastFunc CastFunc
}

func NewVariable(name, value, cast string) (v *Variable, err error) {
	v = &Variable{
		Name: name,
		Value: value,
	}
	switch cast {
	case VariableCastString:
		v.CastFunc = CastString
	case VariableCastJson:
		v.CastFunc = CastJson
	}
	return
}

func (v *Variable) Render() (r string, err error) {
	var value string
	if value, err = v.CastFunc(v.Value); err != nil {
		return
	}
	r = fmt.Sprintf("%s=%s", v.Name, value)
	return
}

type CastFunc func(in string) (r string, err error)

func CastString(in string) (r string, err error) {
	r = fmt.Sprintf("%s", in)
	return
}

func CastJson(in string) (r string, err error) {
	var tmp interface{}
	if err = json.Unmarshal([]byte(in), &tmp); err != nil {
		return
	}
	var data []byte
	if data, err = json.Marshal(&tmp); err != nil {
		return
	}
	r = fmt.Sprintf("'%s'", string(data))
	return
}
