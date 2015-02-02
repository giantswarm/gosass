package compiler

import (
    "os"
    "fmt"
    "log"
)

func fileLog(bail bool, filename string, message string, args ...interface{}) {
    log.Printf(fmt.Sprintf("[%s] %s", filename, message), args...)

    if bail {
        os.Exit(1)
    }
}

func fileLogCompilationError(filename string, err error) {
    fileLog(false, filename, "Could not compile: %s", err.Error())
}

func fileLogCompilation(fromPath string, toPath string) {
    fileLog(false, fromPath, "Compiled to %s", toPath)
}
