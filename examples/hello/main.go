package main

import (
	"fmt"

	"github.com/maximekuhn/pargs"
)

type args struct {
	Friend bool `pargs:"flag:friend"`
}

func main() {
	var args args
	if err := pargs.Parse(&args, nil); err != nil {
		panic(err)
	}

	if args.Friend {
		fmt.Println("Hello, friend !")
	} else {
		fmt.Println("Hello stranger !")
	}
}
