package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Args struct {
	Verbose bool `pargs:"verbose"`
	Help    bool `pargs:"help"`
}

func (a Args) print() {
	fmt.Printf("a.Verbose: %v\n", a.Verbose)
	fmt.Printf("a.Help: %v\n", a.Help)
}

func main() {
	var args Args
	if err := Parse(&args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	args.print()
}

func Parse(target any) error {
	return parse(os.Args[1:], target)
}

func parse(osArgs []string, target any) error {
	targetType := reflect.TypeOf(target).Elem()
	fields := fields(targetType)
	marked := make([]bool, len(osArgs))
	for _, field := range fields {
		for i, arg := range osArgs {
			// TODO: support short flags
			cArg, found := strings.CutPrefix(arg, "--")
			if !found {
				cArg = arg
			}

			if field.Tag != cArg {
				continue
			}

			if marked[i] {
				return fmt.Errorf("arg at index %d is already marked (duplicate?)", i)
			}

			if field.Type != reflect.Bool.String() {
				return fmt.Errorf("TODO: handle other types than boolean")
			}

			val := reflect.ValueOf(target).Elem().FieldByName(field.Name)
			val.SetBool(true)

			marked[i] = true
		}
	}

	totalMarked := 0
	for _, m := range marked {
		if m {
			totalMarked++
		}
	}
	if totalMarked != len(osArgs) {
		return errors.New("not all args have been recognized")
	}

	return nil
}

func fields(t reflect.Type) []Field {
	tags := make([]Field, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("pargs")
		tags = append(tags, Field{
			Type: field.Type.Name(),
			Name: field.Name,
			Tag:  tag,
		})
	}
	return tags
}

type Field struct {
	Type string
	Name string
	Tag  string
}

func (f Field) print() {
	fmt.Printf("f.Type: %v\n", f.Type)
	fmt.Printf("f.Name: %v\n", f.Name)
	fmt.Printf("f.Tag: %v\n", f.Tag)
}
