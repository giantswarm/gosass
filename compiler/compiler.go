package compiler

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	MAX_CONCURRENT_COMPILES = 4
)

// Finds what files are sass compilable in the context's `inputPath`.
func findCompilable(ctx *SassContext) map[string]string {
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

func captureStdout(wg sync.WaitGroup, errorChan chan error, resultsChan chan string, stdout io.ReadCloser) {
	wg.Add(1)
	defer wg.Done()

	stdoutBytes, err := ioutil.ReadAll(stdout)

	if err != nil {
		errorChan <- err
	}

	resultsChan <- string(stdoutBytes)
}

func captureStderr(wg sync.WaitGroup, errorChan chan error, stderr io.ReadCloser) {
	wg.Add(1)
	defer wg.Done()

	stderrScanner := bufio.NewScanner(stderr)

	for stderrScanner.Scan() {
		log.Print(stderrScanner.Text())
	}

	if err := stderrScanner.Err(); err != nil {
		errorChan <- err
	}
}

// Compiles an individual file
func compile(ctx *SassContext, inputPath string, outputPath string) error {
	// Create the parent directory
	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)

	if err != nil {
		return err
	}

	// Create the command and grab stdout/stderr
	cmd := ctx.cmd.Create(inputPath)

	// Grab stdout
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	defer stdout.Close()

	// Grab stderr
	stderr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	defer stderr.Close()

	// Run the command
	err = cmd.Start()

	if err != nil {
		return err
	}

	// Handle stdout/stderr
	var wg sync.WaitGroup
	errorChan := make(chan error, 2)
	stdoutChan := make(chan string, 1)
	go captureStdout(wg, errorChan, stdoutChan, stdout)
	go captureStderr(wg, errorChan, stderr)

	err = cmd.Wait()

	if err != nil {
		return err
	}

	wg.Wait()

	if len(errorChan) > 0 {
		return <-errorChan
	}

	stdoutString := <-stdoutChan

	for _, plugin := range ctx.plugins {
		objs, err := plugin.Objects()

		if err != nil {
			return err
		}

		for _, obj := range objs {
			var newStdoutString string
			err = plugin.Call(fmt.Sprintf("%s.ProcessCss", obj), stdoutString, &newStdoutString)

			if err != nil {
				return err
			}

			stdoutString = newStdoutString
		}
	}

	return ioutil.WriteFile(outputPath, []byte(stdoutString), os.ModePerm)
}

// Compiles many files, as a mapping of input file path -> output file path
func compileMany(ctx *SassContext, mapping map[string]string) bool {
	remaining := len(mapping)
	lock := make(chan bool, MAX_CONCURRENT_COMPILES)
	errorChans := make(map[string]chan error, remaining*2)

	for inputPath, outputPath := range mapping {
		errorChan := make(chan error, 1)
		errorChans[inputPath] = errorChan

		go func(inputPath string, outputPath string, errorChan chan error) {
			lock <- true
			defer func() { <-lock }()

			err := compile(ctx, inputPath, outputPath)
			errorChan <- err
		}(inputPath, outputPath, errorChan)
	}

	hasErrors := false

	for inputPath, errorChan := range errorChans {
		err := <-errorChan

		if err != nil {
			fileLogCompilationError(inputPath, err)
			hasErrors = true
		} else {
			fileLogCompilation(inputPath, mapping[inputPath])
		}
	}

	return hasErrors
}
