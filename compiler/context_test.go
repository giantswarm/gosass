package compiler

import (
    "testing"
    "path/filepath"
    "os"
)

func TestRelativeResolveOutputPath(t *testing.T) {
    t.Parallel()

    ctx := NewSassContext(NewSassCommand(), "foo", "bar")

    p := ctx.resolveOutputPath("foo/baz/bez.scss")

    if p != "bar/baz/bez.css" {
        t.Errorf("Unexpected output path: %s", p)
    }

    p = ctx.resolveOutputPath("foo/baz.ext")

    if p != "bar/baz.ext" {
        t.Errorf("Unexpected output path: %s", p)
    }
}

func TestAbsoluteResolveOutputPath(t *testing.T) {
    t.Parallel()

    wd, err := os.Getwd()

    if err != nil {
        t.Error(err)
    }

    ctx := NewSassContext(NewSassCommand(), "foo", "bar")

    p := ctx.resolveOutputPath(filepath.Join(wd, "foo/baz/bez.scss"))

    if p != "bar/baz/bez.css" {
        t.Errorf("Unexpected output path: %s", p)
    }
}
