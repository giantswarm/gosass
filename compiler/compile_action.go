package compiler

import (
    "os"
)

// CLI endpoint for compiling
func Compile(ctx SassContext) {
    if compileMany(ctx, findCompilable(ctx)) {
        os.Exit(1)
    }
}
