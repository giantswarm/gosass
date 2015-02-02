package compiler

import (
    "os"
)

func Compile(ctx SassContext) {
    if compileMany(ctx, findCompilable(ctx)) {
        os.Exit(1)
    }
}
