package compiler

import (
    "path/filepath"
)

func isSassFile(path string) bool {
    ext := filepath.Ext(path)
    return ext == ".scss"
}

func isPrivateFile(path string) bool {
    base := filepath.Base(path)
    return len(base) > 0 && base[0] == '_'
}
