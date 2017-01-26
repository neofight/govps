package main

import (
	"fmt"

	"github.com/neofight/govps/ssh"
)

type publishMVC struct {
	domain string
}

func (step publishMVC) Execute(cxt context) error {

	session, err := cxt.Client.NewSession()

	if err != nil {
		return fmt.Errorf("Unable to create session: %v", err)
	}

	defer session.Close()

	err = ssh.RunSudoCommand(session, "mkdir /var/www/"+step.domain, cxt.password)

	if err != nil {
		return fmt.Errorf("Unable to create directory: %v", err)
	}

	return ssh.ScpUploadAsRoot(cxt.Client, ".", "/var/www/"+step.domain, cxt.password)
}
