package stparser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseStructTag_Ok(t *testing.T) {
	testcases := []struct {
		TestName  string
		StructTag string
		Kind      reflect.Kind
		Expected  Argument
	}{
		{
			TestName:  "Basic",
			StructTag: "flag:name",
			Kind:      reflect.String,
			Expected:  Argument{Flags: []string{"--name"}},
		},
		{
			TestName:  "Short and long flag",
			StructTag: "flag:[name;n]",
			Kind:      reflect.String,
			Expected:  Argument{Flags: []string{"--name", "-n"}},
		},
		{
			TestName:  "Mandatory",
			StructTag: "flag:name,mandatory",
			Kind:      reflect.String,
			Expected: Argument{
				Flags:     []string{"--name"},
				Mandatory: true,
			},
		},
		{
			TestName:  "Default value",
			StructTag: "flag:gopher,default:true",
			Kind:      reflect.Bool,
			Expected: Argument{
				Flags:        []string{"--gopher"},
				Mandatory:    false,
				DefaultValue: true,
			},
		},
		{
			TestName:  "Negative flag",
			StructTag: "flag:sleep,negative",
			Kind:      reflect.Bool,
			Expected:  Argument{Flags: []string{"--sleep", "--no-sleep"}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.TestName, func(t *testing.T) {
			actual, err := ParseStructTag(tc.StructTag, tc.Kind)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_ParseStructTag_Error(t *testing.T) {
	testcases := []struct {
		TestName  string
		StructTag string
		Kind      reflect.Kind
		Expected  string
	}{
		{
			TestName:  "Unknown attribute",
			StructTag: "hello:world",
			Kind:      reflect.String,
			Expected:  "unknown attribute: 'hello:world'",
		},
		{
			TestName:  "Bad default value type",
			StructTag: "default:true",
			Kind:      reflect.Int, // no compatible with 'true'
			Expected:  "strconv.Atoi: parsing \"true\": invalid syntax",
		},
		{
			TestName:  "Wrong flag attribute syntax",
			StructTag: "flag:[name;n",
			Kind:      reflect.String,
			Expected:  "flag attribute: missing closing bracket",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.TestName, func(t *testing.T) {
			_, err := ParseStructTag(tc.StructTag, tc.Kind)
			assert.Error(t, err)
			assert.Equal(t, tc.Expected, err.Error())
		})
	}
}
