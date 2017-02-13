package tasks

import (
	"fmt"
)

type InstallMono struct {
}

func (task InstallMono) Execute(cxt Context) error {

	err := uploadTemplate(cxt, "systemd", monoUnitFile, cxt.Domain, monoUnitFilePath)

	if err != nil {
		return fmt.Errorf("Failed to deploy Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service unit file uploaded")

	_, err = cxt.VPS.RunSudoCommands(enableMonoServiceCommand)

	if err != nil {
		return fmt.Errorf("Unable to enable Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service enabled")

	_, err = cxt.VPS.RunSudoCommands(startMonoServiceCommand)

	if err != nil {
		return fmt.Errorf("Unable to start Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service started")

	return nil
}

const monoUnitFile = `[Unit]
Description=Mono FastCGI Server

[Service]
Type=simple
ExecStart=/usr/bin/fastcgi-mono-server4 --appconfigfile=/etc/xsp4/debian.webapp /socket=tcp:127.0.0.1:9000

[Install]
WantedBy=multi-user.target`

const monoUnitFilePath = "/lib/systemd/system/mono-fastcgi.service"

const enableMonoServiceCommand = "systemctl enable mono-fastcgi"

const startMonoServiceCommand = "systemctl start mono-fastcgi"
