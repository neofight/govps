package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/neofight/govps/ssh"
)

type installMonoFastCGIService struct {
}

func (step installMonoFastCGIService) Execute(cxt context) error {

	template, err := template.New("systemd").Parse(systemdTemplate)

	if err != nil {
		return fmt.Errorf("Unable to parse systemd template: %v", err)
	}

	var buffer bytes.Buffer

	template.Execute(&buffer, cxt.domain)

	err = ssh.ScpUploadDataAsRoot(cxt.Client, buffer.String(), "/lib/systemd/system/mono-fastcgi.service", cxt.password)

	if err != nil {
		return fmt.Errorf("Failed to upload file: %v", err)
	}

	session, err := cxt.Client.NewSession()

	if err != nil {
		return fmt.Errorf("Unable to create session: %v", err)
	}

	defer session.Close()

	err = ssh.RunSudoCommand(session, "sudo systemctl enable mono-fastcgi", cxt.password)

	if err != nil {
		return fmt.Errorf("Unable to enable Mono FastCGI service: %v", err)
	}

	session, err = cxt.Client.NewSession()

	if err != nil {
		return fmt.Errorf("Unable to create session: %v", err)
	}

	defer session.Close()

	err = ssh.RunSudoCommand(session, "sudo systemctl start mono-fastcgi", cxt.password)

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
