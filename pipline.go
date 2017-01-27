package main

import (
	"golang.org/x/crypto/ssh"
)

type context struct {
	Client   *ssh.Client
	password []byte
	domain   string
}

type step interface {
	Execute(cxt context) error
}

func executePipeline(cxt context, steps []step) error {

	for _, step := range steps {
		err := step.Execute(cxt)

		if err != nil {
			return err
		}
	}

	return nil
}
