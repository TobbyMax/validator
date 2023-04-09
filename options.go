package validator

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	Min = "min"
	Max = "max"
	Len = "len"
	In  = "in"
)

type Options struct {
	// Numeric map to store 'min', 'max' and 'len' options
	Numeric map[string]int
	// InStr slice of string values in 'in' option
	// also usable for printing values from 'in' values in case of error
	InStr []string
	// InInt slice of integers, if 'in' option is applied to an integer
	InInt []int
	Field reflect.StructField
}

// getNumericOption parses 'len', 'max', 'min' options
func (o *Options) getNumericOption(opt string, val string) error {
	n, err := strconv.Atoi(val)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	o.Numeric[opt] = n
	return nil
}

// parseInOption creates only slice InStr if 'kind' is a string
// or slices InStr and InInt, if 'kind' is an integer
func (o *Options) parseInOption(kind reflect.Kind, s string) error {
	o.InStr = strings.Split(s, ",")
	if kind == reflect.Int {
		o.InInt = make([]int, len(o.InStr))
		for i, v := range o.InStr {
			n, err := strconv.Atoi(v)
			if err != nil {
				return ErrInvalidValidatorSyntax
			}
			o.InInt[i] = n
		}
	}
	return nil
}

func (o *Options) setOption(kind reflect.Kind, opt string, val string) error {
	numeric := []string{Min, Max, Len}
	if len(val) == 0 {
		return ErrInvalidValidatorSyntax
	}
	switch {
	case contains(numeric, opt):
		err := o.getNumericOption(opt, val)
		if err != nil {
			return err
		}
	case opt == In:
		err := o.parseInOption(kind, val)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseOptions parses tag string to retrieve constraint options
func parseOptions(kind reflect.Kind, st string) (Options, error) {
	opts := Options{Numeric: make(map[string]int)}
	for _, s := range strings.Split(st, ";") {
		optVal := strings.Split(s, ":")
		if len(optVal) != 2 {
			return Options{}, ErrInvalidValidatorSyntax
		}
		opt, val := strings.Trim(optVal[0], " "), strings.Trim(optVal[1], " ")
		err := opts.setOption(kind, opt, val)
		if err != nil {
			return Options{}, err
		}
	}
	return opts, nil
}

func contains[T comparable](set []T, val T) bool {
	for _, v := range set {
		if v == val {
			return true
		}
	}
	return false
}
