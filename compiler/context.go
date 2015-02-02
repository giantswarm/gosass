package compiler

import (
    "path/filepath"
)

// Stores contextual information for CLI invocations
type SassContext struct {
    cmd *SassCommand
    inputPath string
    outputPath string
}

func NewSassContext(cmd *SassCommand, inputPath string, outputPath string) SassContext {
    return SassContext {
        cmd: cmd,
        inputPath: inputPath,
        outputPath: outputPath,
    }
}

// Gets the equivalent output path for the given path. The given path must be
// within the `inputPath`, but it may be in absolute or relative form.
func (self SassContext) resolveOutputPath(p string) string {
    if filepath.IsAbs(p) {
        absInput, err := filepath.Abs(self.inputPath)

        if err != nil {
            panic(err)
        }

        p, err = filepath.Rel(absInput, p)

        if err != nil {
            panic(err)
        }
    } else {
        np, err := filepath.Rel(self.inputPath, p)

        if err != nil {
            panic(err)
        }

        p = np //QED
    }

    // Replace .scss with .css
    np := filepath.Join(self.outputPath, p)
    ext := filepath.Ext(np)

    if ext == ".scss" {
        np = np[0:len(np)-len(ext)] + ".css"
    }

    return np
}
