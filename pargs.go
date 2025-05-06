package pargs

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
)

var (
	ErrInvalidTargetType    = errors.New("target must be a struct to a pointer")
	ErrStructFieldNotTagged = errors.New("struct field not tagged")
)

type Opts struct {
	Input []string
}

func Parse(target any, o *Opts) error {
	var args []string
	if o != nil && o.Input != nil && len(o.Input) > 0 {
		args = o.Input
	} else {
		args = os.Args[1:]
	}
	return parse(target, args)
}

func parse(target any, args []string) error {
	tfs, err := targetFields(target)
	if err != nil {
		return err
	}

	processed := make(map[string]struct{})

	out := reflect.ValueOf(target).Elem()

	for _, tf := range tfs {
		for _, arg := range args {
			if _, found := processed[arg]; found {
				continue
			}

			if !slices.Contains(tf.flags, arg) {
				continue
			}

			kind := tf.field.Type.Kind()
			switch kind {
			case reflect.Bool:
				f := out.FieldByName(tf.field.Name)
				f.SetBool(true)
			default:
				panic(fmt.Sprintf("TODO: handle kind '%s'", kind))
			}

			processed[arg] = struct{}{}
		}
	}

	return nil
}

func targetFields(target any) ([]targetField, error) {
	if reflect.TypeOf(target).Kind() != reflect.Pointer {
		return nil, ErrInvalidTargetType
	}
	if reflect.ValueOf(target).Elem().Kind() != reflect.Struct {
		return nil, ErrInvalidTargetType
	}

	t := reflect.TypeOf(target).Elem()
	targetFields := make([]targetField, 0)
	for i := range t.NumField() {
		field := t.Field(i)
		targetField, err := parseTargetField(field)
		if err != nil {
			return nil, err
		}
		targetFields = append(targetFields, targetField)
	}
	return targetFields, nil
}

func parseTargetField(field reflect.StructField) (targetField, error) {
	tf := targetField{field: field}

	pargsTag := field.Tag.Get("pargs")
	if pargsTag == "" {
		return tf, ErrStructFieldNotTagged
	}

	parts := strings.Split(pargsTag, ",")

	for _, part := range parts {
		if strings.HasPrefix(part, "flag:") {
			flags, err := parseFlags(part)
			if err != nil {
				return tf, err
			}
			tf.flags = flags
		}

		if part == "mandatory" {
			tf.mandatory = true
		}
	}

	return tf, nil
}

func parseFlags(s string) ([]string, error) {
	const (
		prefix = "flag:"
		sep    = ";"
	)
	after, found := strings.CutPrefix(s, prefix)
	if !found {
		panic("prefix not found")
	}

	if len(after) == 0 {
		return nil, errors.New("missing flag value")
	}

	flags := make([]string, 0)
	parts := strings.Split(after, sep)
	for _, part := range parts {
		if len(part) == 1 {
			flags = append(flags, fmt.Sprintf("-%s", part))
		} else {
			flags = append(flags, fmt.Sprintf("--%s", part))
		}
	}

	return flags, nil
}

type targetField struct {
	flags     []string
	mandatory bool
	field     reflect.StructField
}
