package autooverload

import "github.com/aybolid/envar"

func init() {
	if err := envar.Overload(); err != nil {
		panic(err)
	}
}
