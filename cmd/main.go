package main

import (
	"fmt"

	"github.com/aybolid/envar"
)

type ENV struct {
	FOO    string
	BAR    string
	NUMBER int
	LIST   []string
	N_LIST []int
}

func main() {
	envFile := "test/fixtures/loadin.env"

	env := new(ENV)
	err := envar.LoadIn(envFile, env)
	if err != nil {
		panic(err)
	}
	fmt.Println(env)
}
