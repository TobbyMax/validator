package validator

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	s := make([]string, len(v))
	for i, err := range v {
		s[i] = err.Err.Error()
	}
	return strings.Join(s, ";\n")
}

func Validate(v any) error {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	vld := Validator{}
	if t.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	for i := 0; i < t.NumField(); i++ {
		vld.ValidateField(t.Field(i), val.Field(i))
	}
	if len(vld.Errors) != 0 {
		return vld.Errors
	}
	return nil
}
