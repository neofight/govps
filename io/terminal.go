package io

import (
	sshTerminal "golang.org/x/crypto/ssh/terminal"
)

type terminal interface {
	ReadPassword(fd int) ([]byte, error)
}

type realTerminal struct {
}

func (realTerminal) ReadPassword(fd int) ([]byte, error) {
	return sshTerminal.ReadPassword(fd)
}

var Terminal terminal = realTerminal{}
