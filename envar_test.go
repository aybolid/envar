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
		if v != expectedValue {
			t.Errorf("missmatch for %q key: expected %q, got %q\n", k, expectedValue, v)
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
