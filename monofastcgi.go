package main

import (
	"fmt"

	"github.com/neofight/govps/ssh"
)

type installMonoFastCGIService struct {
}

func (step installMonoFastCGIService) Execute(cxt context) error {

	err := uploadTemplate(cxt, "systemd", systemdTemplate, cxt.domain, "/lib/systemd/system/mono-fastcgi.service")

	if err != nil {
		return fmt.Errorf("Failed to deploy Mono FastCGI service: %v", cxt.domain, err)
	}

	err = ssh.RunSudoCommand(cxt.Client, "sudo systemctl enable mono-fastcgi", cxt.password)

	if err != nil {
		return fmt.Errorf("Unable to enable Mono FastCGI service: %v", err)
	}

	err = ssh.RunSudoCommand(cxt.Client, "sudo systemctl start mono-fastcgi", cxt.password)

	if err != nil {
		return fmt.Errorf("Unable to start Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service created")

	return nil
}

var systemdTemplate = `[Unit]
Description=Mono FastCGI Server

[Service]
Type=simple
ExecStart=/usr/bin/fastcgi-mono-server4 /applications=/:/var/www/{{.}}/ /socket=tcp:127.0.0.1:9000

[Install]
WantedBy=multi-user.target`
