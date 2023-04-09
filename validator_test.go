package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrInvalidValidatorSyntax.Error()
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int    `validate:"in:20,25,30"`
					InNeg     int    `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int    `validate:"min:10"`
					MinIntNeg int    `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int    `validate:"max:20"`
					MaxIntNeg int    `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int    `validate:"in:-1,-3,5,7"`
					InD     int    `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int    `validate:"min:-12"`
					MinC int    `validate:"min:5-"`
					MinD int    `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int    `validate:"max:-12"`
					MaxD int    `validate:"max:5-"`
					MaxE int    `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
		{
			name: "slices: valid struct with tagged",
			args: args{
				v: struct {
					Len       []string `validate:"len:5"`
					LenZ      []string `validate:"len:0"`
					InInt     []int    `validate:"in:20,25,30"`
					InNeg     []int    `validate:"in:-20,-25,-30"`
					InStr     []string `validate:"in:foo,bar"`
					MinInt    []int    `validate:"min:10"`
					MinIntNeg []int    `validate:"min:-10"`
					MinStr    []string `validate:"min:4"`
					MinStrNeg []string `validate:"min:-1"`
					MaxInt    []int    `validate:"max:20"`
					MaxIntNeg []int    `validate:"max:-2"`
					MaxStr    []string `validate:"max:20"`
				}{
					Len:       []string{"abcde", "12345", "yvtro", "bussy"},
					LenZ:      []string{"", ""},
					InInt:     []int{25, 20, 25},
					InNeg:     []int{-25, -30, -20, -20},
					InStr:     []string{"bar", "foo", "foo"},
					MinInt:    []int{25, 20, 11, 33, 100},
					MinIntNeg: []int{0, 2, -5, -2},
					MinStr:    []string{"abcde", "1234", "yvt3134rro", "bussy30"},
					MinStrNeg: []string{"abc", "", "dkjsfh"},
					MaxInt:    []int{-25, -20, 11, -33, 13},
					MaxIntNeg: []int{-25, -20, -11, -33, -100},
					MaxStr:    []string{"abcde12345678990", "1234", "yvt3134rro", "bussy330"},
				},
			},
			wantErr: false,
		},
		{
			name: "slices: wrong length",
			args: args{
				v: struct {
					Lower    []string `validate:"len:24"`
					Higher   []string `validate:"len:5"`
					Zero     []string `validate:"len:3"`
					BadSpec  []string `validate:"len:%12"`
					Negative []string `validate:"len:-6"`
				}{
					Lower:    []string{"abcdef", "dwayne", "1234567890qwertyuiopasdf"},
					Higher:   []string{"abcdef", "rock", "o", "kanye"},
					Zero:     []string{"", "abs", "b2b"},
					BadSpec:  []string{"abc", "kj"},
					Negative: []string{"abcd"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 8)
				return true
			},
		},
		{
			name: "slices: wrong in",
			args: args{
				v: struct {
					InA     []string `validate:"in:ab,cd"`
					InB     []string `validate:"in:aa,bb,cd,ee"`
					InC     []int    `validate:"in:-1,-3,5,7"`
					InD     []int    `validate:"in:5-"`
					InEmpty []string `validate:"in:"`
				}{
					InA:     []string{"ef", "ab", "cd", "hh"},
					InB:     []string{"ab", "aa", "ye"},
					InC:     []int{2, -3, 7, 7, -2},
					InD:     []int{12, 1, 1, 1, 1, 1, 12},
					InEmpty: []string{""},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 8)
				return true
			},
		},
		{
			name: "slices: wrong min",
			args: args{
				v: struct {
					MinA []string `validate:"min:12"`
					MinB []int    `validate:"min:-12"`
					MinC []int    `validate:"min:5-"`
					MinD []int    `validate:"min:"`
					MinE []string `validate:"min:"`
				}{
					MinA: []string{"ef", "jesus and judas", "human"},
					MinB: []int{-22, -38, -11},
					MinC: []int{12, 17, 1},
					MinD: []int{11, 13, 12},
					MinE: []string{"abc"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 7)
				return true
			},
		},
		{
			name: "slices: wrong max",
			args: args{
				v: struct {
					MaxA []string `validate:"max:2"`
					MaxB []string `validate:"max:-7"`
					MaxC []int    `validate:"max:-12"`
					MaxD []int    `validate:"max:5-"`
					MaxE []int    `validate:"max:"`
					MaxF []string `validate:"max:"`
				}{
					MaxA: []string{"efgh", "17", "777"},
					MaxB: []string{"ab", ""},
					MaxC: []int{22, 11, -33, -11},
					MaxD: []int{12, 12},
					MaxE: []int{11, 11, 11},
					MaxF: []string{"abc"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 10)
				return true
			},
		},
		{
			name: "valid struct with untagged fields: not string or int",
			args: args{
				v: struct {
					f1 float32
					f2 float32
				}{},
			},
			wantErr: false,
		},
		{
			name: "unable to validate type",
			args: args{
				v: struct {
					F float64 `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == "field of type float64 can not be validated"
			},
		},
		{
			name: "multiple constraints: valid",
			args: args{
				v: struct {
					MinMax int `validate:"min:12; max: 15"`
				}{
					MinMax: 13,
				},
			},
			wantErr: false,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 1)
				return true
			},
		},
		{
			name: "multiple constraints: invalid",
			args: args{
				v: struct {
					MinMax string `validate:"min:12; max:15"`
				}{
					MinMax: "spkjishu",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 1)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
