package validator

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

var ErrFieldNotValid = "field '%s' of type %s is not valid: has constraint ('%s': %v), but got value = %v"

type FieldValidator interface {
	ValidateField(sf reflect.StructField, val reflect.Value)
}

type Validator struct {
	Errors ValidationErrors
}

func (v *Validator) ValidateField(sf reflect.StructField, val reflect.Value) {
	if sf.Tag.Get("validate") == "" {
		return
	}
	if !sf.IsExported() {
		v.Errors = append(v.Errors, ValidationError{ErrValidateForUnexportedFields})
		return
	}
	if sf.Type.Kind() == reflect.Slice {
		v.validateSlice(sf, val)
		return
	}
	opts := v.getFieldOptions(sf.Type.Kind(), sf.Tag.Get("validate"))
	v.validateValue(val, sf, opts)
}

// getFieldOptions returns options for string and int fields (or slices),
// appends an error to Validator.Errors in case of ErrInvalidValidatorSyntax
// and fields of types, which differ from string and int, being tagged
func (v *Validator) getFieldOptions(kind reflect.Kind, tag string) Options {
	var (
		opts Options
		err  error
	)
	if kind == reflect.String || kind == reflect.Int {
		opts, err = parseOptions(kind, tag)
	} else {
		err = errors.Errorf("field of type %s can not be validated", kind)
	}
	if err != nil {
		v.Errors = append(v.Errors, ValidationError{err})
	}
	return opts
}

// validateValue validates a value according to given options
func (v *Validator) validateValue(val reflect.Value, sf reflect.StructField, opts Options) {
	v.validateIn(val, sf, opts)
	v.validateNumeric(val, sf, opts)
}

// validateSlice gets options from a tag and validates all values in slice according to these options
func (v *Validator) validateSlice(sf reflect.StructField, sl reflect.Value) {
	opts := v.getFieldOptions(sf.Type.Elem().Kind(), sf.Tag.Get("validate"))
	for i := 0; i < sl.Len(); i++ {
		v.validateValue(sl.Index(i), sf, opts)
	}
}

// validateIn checks if value corresponds to 'in' constraint
func (v *Validator) validateIn(val reflect.Value, sf reflect.StructField, opts Options) {
	switch val.Kind() {
	case reflect.Int:
		if opts.InInt != nil && !contains(opts.InInt, int(val.Int())) {
			v.Errors = append(v.Errors, ValidationError{
				errors.Errorf(ErrFieldNotValid, sf.Name, sf.Type.Kind(), In,
					strings.Join(opts.InStr, ","), val.Int())})
		}
	case reflect.String:
		if opts.InStr != nil && !contains(opts.InStr, val.String()) {
			v.Errors = append(v.Errors, ValidationError{
				errors.Errorf(ErrFieldNotValid, sf.Name, sf.Type.Kind(), In,
					strings.Join(opts.InStr, ","), val.String())})
		}
	}
}

// validateNumeric checks if value corresponds to 'len', 'min' and 'max' constraints
func (v *Validator) validateNumeric(val reflect.Value, sf reflect.StructField, opts Options) {
	n := v.getNumericValue(val)
	for k, l := range opts.Numeric {
		switch {
		case k == Min && n < l:
			fallthrough
		case k == Max && n > l:
			fallthrough
		case k == Len && n != l:
			v.Errors = append(v.Errors, ValidationError{
				errors.Errorf(ErrFieldNotValid, sf.Name, sf.Type.Kind(), k, l, n)})
		}
	}
}

// getNumericValueAndMessage is a supporting function, which returns underlying value for integers
// and length for strings, and also returns message in case of errors
func (v *Validator) getNumericValue(val reflect.Value) int {
	switch val.Kind() {
	case reflect.Int:
		return int(val.Int())
	case reflect.String:
		return len(val.String())
	}
	return 0
}
