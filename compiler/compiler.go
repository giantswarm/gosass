package compiler

import (
    "os"
    "io"
    "log"
    "bufio"
    "path/filepath"
)

// Finds what files are sass compilable in the context's `inputPath`.
func findCompilable(ctx SassContext) map[string]string {
    compilable := make(map[string]string, 100)

    filepath.Walk(ctx.inputPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fileLogCompilationError(path, err)
        } else if !info.IsDir() && isSassFile(path) && !isPrivateFile(path) {
            compilable[path] = ctx.resolveOutputPath(path)
        }

        return nil
    })

    return compilable
}

// Handles stdin/stderr from a process and pipes it to the log
func redirectProcessOutput(name string, pipe io.ReadCloser) {
    scanner := bufio.NewScanner(pipe)

    for scanner.Scan() {
        log.Printf("[%s] %s", name, scanner.Text())
    }

    err := scanner.Err(); if err != nil {
        log.Printf("Could not read from %s: %s", name, err.Error())
    }
}

// Compiles an individual file
func compile(ctx SassContext, inputPath string, outputPath string) error {
    err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)

    if err != nil {
        return err
    }

    cmd := ctx.cmd.Create(inputPath, outputPath)

    stdout, err := cmd.StdoutPipe()

    if err != nil {
        return err
    }

    stderr, err := cmd.StderrPipe()

    if err != nil {
        return err
    }

    err = cmd.Start()

    if err != nil {
        return err
    }

    redirectProcessOutput("stdout", stdout)
    redirectProcessOutput("stderr", stderr)
    return cmd.Wait()
}

// Compiles many files, as a mapping of input file path -> output file path
func compileMany(ctx SassContext, mapping map[string]string) bool {
    remaining := len(mapping)
    errorChans := make(map[string]chan error, remaining * 2)

    for inputPath, outputPath := range mapping {
        errorChan := make(chan error, 1)
        errorChans[inputPath] = errorChan

        go func(inputPath string, outputPath string, errorChan chan error) {
            errorChan <- compile(ctx, inputPath, outputPath)
        }(inputPath, outputPath, errorChan)
    }

    hasErrors := false

    for inputPath, errorChan := range errorChans {
        err := <- errorChan

        if err != nil {
            fileLogCompilationError(inputPath, err)
            hasErrors = true
        } else {
            fileLogCompilation(inputPath, mapping[inputPath])
        }
    }

    return hasErrors
}
