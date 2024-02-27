package autoload

import "github.com/aybolid/envar"

func init() {
	if err := envar.Load(); err != nil {
		panic(err)
	}
}
