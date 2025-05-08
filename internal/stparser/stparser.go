package stparser

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Argument struct {
	Flags        []string
	Mandatory    bool
	DefaultValue any
}

func ParseStructTag(st string, kind reflect.Kind) (Argument, error) {
	arg := Argument{}
	negative := false
	parts := strings.Split(st, ",")
	for _, part := range parts {
		if strings.HasPrefix(part, "flag:") {
			flags, err := parseFlagValue(part[len("flag:"):])
			if err != nil {
				return arg, err
			}
			arg.Flags = flags
			continue
		}

		if strings.HasPrefix(part, "default:") {
			val, err := parseDefaultValue(part[len("default:"):], kind)
			if err != nil {
				return arg, err
			}
			arg.DefaultValue = val
			continue
		}

		if part == "mandatory" {
			arg.Mandatory = true
			continue
		}

		if part == "negative" {
			negative = true
			continue
		}

		return arg, fmt.Errorf("unknown attribute: '%s'", part)
	}

	if negative {
		if len(arg.Flags) < 1 {
			return arg, errors.New("no flag attribute provided")
		}
		arg.Flags = append(
			arg.Flags,
			fmt.Sprintf("--no-%s", strings.TrimLeft(arg.Flags[0], "-")),
		)
	}

	return arg, nil
}

func parseFlagValue(flag string) ([]string, error) {
	parts := make([]string, 0)
	if strings.HasPrefix(flag, "[") {
		if !strings.HasSuffix(flag, "]") {
			return nil, errors.New("flag attribute: missing closing bracket")
		}
		withoutBrackets := flag[len("[") : len(flag)-1]
		parts = strings.Split(withoutBrackets, ";")
	} else {
		parts = strings.Split(flag, ";")
	}

	flags := make([]string, 0)
	for _, part := range parts {
		if len(part) == 1 {
			flags = append(flags, fmt.Sprintf("-%s", part))
		} else {
			flags = append(flags, fmt.Sprintf("--%s", part))
		}
	}

	return flags, nil
}

func parseDefaultValue(s string, k reflect.Kind) (any, error) {
	if strings.HasPrefix("kind", k.String()) {
		return nil, errors.New("pargs only supports standard go types for now")
	}

	switch k {
	case reflect.String:
		return s, nil
	case reflect.Bool:
		return strconv.ParseBool(s)
	case reflect.Int:
		return strconv.Atoi(s)
	default:
		return nil, fmt.Errorf("unsupported kind: %s", k)
	}
}
