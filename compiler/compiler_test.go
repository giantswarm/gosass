package compiler

import (
	"io/ioutil"
	"testing"
)

const EXPECTED_SIMPLE_COMPILE = `body {
  font-weight: bold; }
`

func TestFindCompilable(t *testing.T) {
	t.Parallel()

	ctx := NewSassContext(NewSassCommand(), "../integration/bad-src", "../integration/out")
	compilable := findCompilable(ctx)

	if len(compilable) != 0 {
		t.Error()
	}

	ctx = NewSassContext(NewSassCommand(), "../integration/src", "../integration/out")
	compilable = findCompilable(ctx)

	if len(compilable) != 3 {
		t.Error()
	}

	if compilable["../integration/src/02.simple-import.scss"] != "../integration/out/02.simple-import.css" {
		t.Error()
	}

	if compilable["../integration/src/03.multiple-imports.scss"] != "../integration/out/03.multiple-imports.css" {
		t.Error()
	}

	if compilable["../integration/src/01.simple.scss"] != "../integration/out/01.simple.css" {
		t.Error()
	}
}

func TestCompile(t *testing.T) {
	t.Parallel()

	ctx := NewSassContext(NewSassCommand(), "../integration/src", "../integration/out")
	err := compile(ctx, "../integration/src/01.simple.scss", "../integration/out/01.simple.css")

	if err != nil {
		t.Error(err)
	}

	b, err := ioutil.ReadFile("../integration/out/01.simple.css")

	if err != nil {
		t.Error(err)
	}

	if string(b) != EXPECTED_SIMPLE_COMPILE {
		t.Errorf("Unexpected compiled results: %s", string(b))
	}
}
