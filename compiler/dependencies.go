package compiler

import (
    "fmt"
    "regexp"
    "path/filepath"
)

var importPattern = regexp.MustCompile("\\@import (\\'|\\\")([^'\"]+)(\\'|\\\")")

type SassDependencyResolver struct {
    filecache *FileCache
    shallowDeps map[string][]string
    deepDeps map[string][]string
}

func NewSassDependencyResolver(filecache *FileCache) *SassDependencyResolver {
    return &SassDependencyResolver {
        filecache: filecache,
        shallowDeps: make(map[string][]string, 100),
        deepDeps: make(map[string][]string, 100),
    }
}

// Gets the files imported directly by the given file
func (self *SassDependencyResolver) shallowResolve(path string) ([]string, error) {
    deps, ok := self.shallowDeps[path]

    if ok {
        return deps, nil
    }

    // Get the file contents
    contents, err := self.filecache.Get(path)

    if err != nil {
        return nil, err
    }

    // Build the matches
    matches := importPattern.FindAllSubmatch(contents, -1)
    deps = make([]string, len(matches))

    for i, match := range matches {
        ref := string(match[2])
        refPath := filepath.Join(filepath.Dir(ref), fmt.Sprintf("_%s.scss", filepath.Base(ref)))
        deps[i] = filepath.Join(filepath.Dir(path), refPath)
    }

    self.shallowDeps[path] = deps
    return deps, nil
}

// Gets all files imported by the given file, including indirect imports
func (self *SassDependencyResolver) Resolve(path string) ([]string, error) {
    abs, err := filepath.Abs(path)

    if err != nil {
        return nil, err
    }

    deps, ok := self.deepDeps[abs]

    if ok {
        return deps, nil
    }

    scanned := make(map[string]bool, 100)
    unscanned := make(map[string]bool, 100)
    unscanned[abs] = true

    for len(unscanned) > 0 {
        for subpath := range unscanned {
            // Move the file to scanned
            delete(unscanned, subpath)
            scanned[subpath] = true

            // Get the dependencies or read them if needed
            deps, err := self.shallowResolve(subpath)

            if err != nil {
                return nil, err
            }

            // Add the dependency to unscanned if it hasn't been scanned
            // already
            for _, dep := range deps {
                _, ok = scanned[dep]

                if !ok {
                    unscanned[dep] = true
                }
            }
        }
    }

    deps = make([]string, 0, len(scanned))

    for dep := range scanned {
        if dep != abs {
            deps = append(deps, dep)
        }
    }

    self.deepDeps[abs] = deps
    return deps, nil
}

// Gets what files are dependent on the given file, including indirectly
func (self *SassDependencyResolver) ReverseResolve(path string) ([]string, error) {
    abs, err := filepath.Abs(path)

    if err != nil {
        return nil, err
    }

    reverseDeps := make([]string, 0)

    for otherPath, deps := range self.deepDeps {
        for _, dep := range deps {
            if dep == abs {
                reverseDeps = append(reverseDeps, otherPath)
                break
            }
        }
    }

    return reverseDeps, nil
}

// Invalidates the cached entry for the given file
func (self *SassDependencyResolver) Invalidate(path string) error {
    abs, err := filepath.Abs(path)

    if err != nil {
        return err
    }

    delete(self.shallowDeps, abs)
    delete(self.deepDeps, abs)
    return nil
}
