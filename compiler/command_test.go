package compiler

import (
	"testing"
)

const EXPECTED_STDOUT = `Usage: sassc [options] [INPUT] [OUTPUT]

Options:
   -s, --stdin             Read input from standard input instead of an input file.
   -t, --style NAME        Output style. Can be: nested, compressed.
   -l, --line-numbers      Emit comments showing original line numbers.
       --line-comments
   -I, --load-path PATH    Set Sass import path.
   -m, --sourcemap         Emit source map.
   -M, --omit-map-comment  Omits the source map url comment.
   -p, --precision         Set the precision for numbers.
   -v, --version           Display compiled versions.
   -h, --help              Display this help message.

`

// Both tests that the SassCommand struct works, and that the correct sassc
// version is installed
func TestSassCommand(t *testing.T) {
	t.Parallel()

	cmd := NewSassCommand()
	cmd.AddArgument("--help")
	proc := cmd.Create("a")

	stdout, err := proc.Output()

	if err != nil {
		t.Error(err)
	}

	if string(stdout) != EXPECTED_STDOUT {
		t.Errorf("Unexpected stdout: %s", stdout)
	}
}
