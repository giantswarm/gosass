package compiler

import (
    "testing"
)

const EXPECTED_FILE = `body {
    font-weight: bold;
}
`

func TestFileCacheInvalidate(t *testing.T) {
    t.Parallel()

    fc := NewFileCache()

    // Make sure we can arbitrarily invalidate
    err := fc.Invalidate("some/random/path")

    if err != nil {
        t.Error(err)
    }
}

func TestFileCacheGet(t *testing.T) {
    t.Parallel()

    fc := NewFileCache()

    checkFile := func() {
        contents, err := fc.Get("../integration/src/01.simple.scss")

        if err != nil {
            t.Error(err)
        }

        if string(contents) != EXPECTED_FILE {
            t.Errorf("Unexpected file contents: %s", string(contents))
        }
    }

    // Get a file (uncached)
    checkFile()

    // Reget the file - should hit the cache
    checkFile()

    // Invalidate the file
    err := fc.Invalidate("../integration/src/01.simple.scss")

    if err != nil {
        t.Error(err)
    }

    // Get the file one last time - should be uncached again
    checkFile()
}

