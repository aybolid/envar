package envar

import (
	"fmt"
	"testing"
)

type LoadInENV struct {
	FOO    string
	BAR    string
	NUMBER int
	LIST   []string
	N_LIST []int
}

func TestLoadIn(t *testing.T) {
	envFile := "test/fixtures/loadin.env"
	env := new(LoadInENV)

	err := LoadIn(envFile, env)
	if err != nil {
		logError(t, err.Error())
	}

	expectedFOO := "bar"
	if env.FOO != expectedFOO {
		msg := fmt.Sprintf("invalid value for \"FOO\" key. expected: %q; got: %q", expectedFOO, env.FOO)
		logError(t, msg)
	}

	expectedBAR := "foo"
	if env.BAR != expectedBAR {
		msg := fmt.Sprintf("invalid value for \"BAR\" key. expected: %q; got: %q", expectedBAR, env.BAR)
		logError(t, msg)
	}

	expectedNUMBER := 123
	if env.NUMBER != expectedNUMBER {
		msg := fmt.Sprintf("invalid value for \"NUMBER\" key. expected: %d; got: %d", expectedNUMBER, env.NUMBER)
		logError(t, msg)
	}

	expectedLIST := []string{"one", "bar", "three"}
	if !equalSlices(env.LIST, expectedLIST) {
		msg := fmt.Sprintf("invalid value for \"LIST\" key. expected: %v; got: %v", expectedLIST, env.LIST)
		logError(t, msg)
	}

	expectedN_LIST := []int{123, 123, 123, 123}
	if !equalSlicesInt(env.N_LIST, expectedN_LIST) {
		msg := fmt.Sprintf("invalid value for \"N_LIST\" key. expected: %v; got: %v", expectedN_LIST, env.N_LIST)
		logError(t, msg)
	}
}

func equalSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func equalSlicesInt(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
