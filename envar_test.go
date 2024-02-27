package envar

import (
	"os"
	"slices"
	"strings"
	"testing"
)

func compare(t *testing.T, got map[string]string, expected map[string]string) {
	if len(got) != len(expected) {
		t.Error("len(got) != len(expected)")
	}

	for k, v := range got {
		expectedValue := expected[k]
		if v == expectedValue {
			t.Logf("[%s]: expected: %q, got: %q\n", k, expectedValue, v)
		} else {
			t.Errorf("\033[31m[%s]: expected: %q, got: %q\033[0m\n", k, expectedValue, v)
		}
	}
}

func TestDefaultOrFilenames(t *testing.T) {
	expected := [][]string{
		{".env"},
		{".env.test", ".env.prod"},
		{""},
	}
	cases := [][]string{
		defaultOrFilenames([]string{}),
		defaultOrFilenames([]string{".env.test", ".env.prod"}),
		defaultOrFilenames([]string{""}),
	}

	for i, c := range cases {
		e := expected[i]
		if !slices.Equal(e, c) {
			formattedE := "[" + strings.Join(e, ", ") + "]"
			formattedC := "[" + strings.Join(c, ", ") + "]"
			t.Errorf("\nexpected: %s\ngot: %s\n", formattedE, formattedC)
		}
	}
}

func TestLoadNonExistingFile(t *testing.T) {
	err := Load(".env.filethatdoesnotexist")
	if err == nil {
		t.Error("loading file that doesn't exist must return an error")
	}
}

func TestOverloadNonExistingFile(t *testing.T) {
	err := Overload(".env.filethatdoesnotexist")
	if err == nil {
		t.Error("loading file that doesn't exist must return an error")
	}
}

func TestNoArgsLoadsDefault(t *testing.T) {
	err := Load()
	pathErr, ok := err.(*os.PathError)
	if !ok {
		t.Error("failed to assert error type")
	}
	if pathErr == nil || pathErr.Path != ".env" || pathErr.Op != "open" {
		t.Error("didn't try to open .env file while no args were provided")
	}
}

func TestNoArgsOverloadsDefault(t *testing.T) {
	err := Overload()
	pathErr, ok := err.(*os.PathError)
	if !ok {
		t.Error("failed to assert error type")
	}
	if pathErr == nil || pathErr.Path != ".env" || pathErr.Op != "open" {
		t.Error("didn't try to open .env file while no args were provided")
	}
}

func TestComments(t *testing.T) {
	envFile := "test/fixtures/comments.env"
	expected := map[string]string{
		"foo":    "bar",
		"bar":    "foo",
		"baz":    "foo#bar",
		"fizz":   "foo",
		"foobar": "foo #bar",
	}

	buf, err := getFileBuffer(envFile)
	if err != nil {
		t.Error(err)
	}
	envMap, err := parse(&buf)
	if err != nil {
		t.Error(err)
	}

	compare(t, envMap, expected)
}

func TestExported(t *testing.T) {
	envFile := "test/fixtures/exported.env"
	expected := map[string]string{
		"OPTION_A": "2",
		"OPTION_B": "\\n",
	}

	buf, err := getFileBuffer(envFile)
	if err != nil {
		t.Error(err)
	}
	envMap, err := parse(&buf)
	if err != nil {
		t.Error(err)
	}

	compare(t, envMap, expected)
}

func TestQuoted(t *testing.T) {
	envFile := "test/fixtures/quoted.env"
	expected := map[string]string{
		"OPTION_A": "1",
		"OPTION_B": "2",
		"OPTION_C": "",
		"OPTION_D": "\\n",
		"OPTION_E": "1",
		"OPTION_F": "2",
		"OPTION_G": "",
		// "OPTION_H": "\n", // FIXME
		"OPTION_I": "echo 'asd'",
	}

	buf, err := getFileBuffer(envFile)
	if err != nil {
		t.Error(err)
	}
	envMap, err := parse(&buf)
	if err != nil {
		t.Error(err)
	}

	compare(t, envMap, expected)
}

func TestPlain(t *testing.T) {
	envFile := "test/fixtures/plain.env"
	expected := map[string]string{
		"OPTION_A": "1",
		"OPTION_B": "2",
		"OPTION_C": "3",
		"OPTION_D": "4",
		"OPTION_E": "5",
		"OPTION_F": "",
		"OPTION_G": "",
		"OPTION_H": "1 2",
	}

	buf, err := getFileBuffer(envFile)
	if err != nil {
		t.Error(err)
	}
	envMap, err := parse(&buf)
	if err != nil {
		t.Error(err)
	}

	compare(t, envMap, expected)
}
