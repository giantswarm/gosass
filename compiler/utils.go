package compiler

import (
	"path/filepath"
)

// Returns whether the given file is sass
func isSassFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".scss"
}

// Returns whether the given file is private (i.e. whether the filename starts
// with an underscore)
func isPrivateFile(path string) bool {
	base := filepath.Base(path)
	return len(base) > 0 && base[0] == '_'
}
