package envar

import (
	"slices"
	"strings"
	"testing"
)

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
		t.Error("loading file that doesn't exist must return error")
	}
}
