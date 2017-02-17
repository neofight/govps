package io

import (
	"io"
	"os"

	sshTerminal "golang.org/x/crypto/ssh/terminal"
)

var StdOut io.Writer = os.Stdout

type terminal interface {
	ReadPassword(fd int) ([]byte, error)
}

type realTerminal struct {
}

func (realTerminal) ReadPassword(fd int) ([]byte, error) {
	return sshTerminal.ReadPassword(fd)
}

var Terminal terminal = realTerminal{}
