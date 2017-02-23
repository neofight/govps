package tasks

import (
	"fmt"
)

type InstallMono struct {
}

func (task InstallMono) Execute(cxt Context) error {

	err := uploadTemplate(cxt, "systemd", MonoUnitFile, cxt.Domain, MonoUnitFilePath)

	if err != nil {
		return fmt.Errorf("Failed to deploy Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service unit file uploaded")

	_, err = cxt.VPS.RunSudoCommands(EnableMonoServiceCommand)

	if err != nil {
		return fmt.Errorf("Unable to enable Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service enabled")

	_, err = cxt.VPS.RunSudoCommands(StartMonoServiceCommand)

	if err != nil {
		return fmt.Errorf("Unable to start Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service started")

	return nil
}

const MonoUnitFile = `[Unit]
Description=Mono FastCGI Server

[Service]
Type=simple
ExecStart=/usr/bin/fastcgi-mono-server4 --appconfigfile=/etc/xsp4/debian.webapp /socket=tcp:127.0.0.1:9000

[Install]
WantedBy=multi-user.target`

const MonoUnitFilePath = "/lib/systemd/system/mono-fastcgi.service"

const EnableMonoServiceCommand = "systemctl enable mono-fastcgi"

const StartMonoServiceCommand = "systemctl start mono-fastcgi"
