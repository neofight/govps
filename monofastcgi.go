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

	fmt.Println("Mono FastCGI service unit file uploaded")

	err = ssh.RunSudoCommands(cxt.Client, cxt.password, "sudo systemctl enable mono-fastcgi")

	if err != nil {
		return fmt.Errorf("Unable to enable Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service enabled")

	err = ssh.RunSudoCommands(cxt.Client, cxt.password, "sudo systemctl start mono-fastcgi")

	if err != nil {
		return fmt.Errorf("Unable to start Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service started")

	return nil
}

var systemdTemplate = `[Unit]
Description=Mono FastCGI Server

[Service]
Type=simple
ExecStart=/usr/bin/fastcgi-mono-server4 /applications=/:/var/www/{{.}}/ /socket=tcp:127.0.0.1:9000

[Install]
WantedBy=multi-user.target`
