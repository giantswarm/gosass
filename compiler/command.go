package compiler

import (
	"os/exec"
)

type SassCommand struct {
	Args []string
}

func NewSassCommand() *SassCommand {
	return &SassCommand{
		Args: make([]string, 0),
	}
}

func (self *SassCommand) AddArgument(arg string) {
	self.Args = append(self.Args, arg)
}

func (self *SassCommand) Create(inputFilePath string) *exec.Cmd {
	args := make([]string, len(self.Args)+1)
	copy(args, self.Args)
	args[len(args)-1] = inputFilePath
	return exec.Command("sassc", args...)
}
