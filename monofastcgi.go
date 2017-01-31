package main

import (
	"encoding/xml"
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

	_, err = ssh.RunSudoCommands(cxt.Client, cxt.password, "sudo systemctl enable mono-fastcgi")

	if err != nil {
		return fmt.Errorf("Unable to enable Mono FastCGI service: %v", err)
	}

	fmt.Println("Mono FastCGI service enabled")

	_, err = ssh.RunSudoCommands(cxt.Client, cxt.password, "sudo systemctl start mono-fastcgi")

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
ExecStart=/usr/bin/fastcgi-mono-server4 --appconfigfile=/etc/xsp4/debian.webapp /socket=tcp:127.0.0.1:9000

[Install]
WantedBy=multi-user.target`

type addSiteToMonoFastCGIConfiguration struct {
}

func (step addSiteToMonoFastCGIConfiguration) Execute(cxt context) error {

	file, err := ssh.ScpDownloadFile(cxt.Client, "/etc/xsp4/debian.webapp")

	if err != nil {
		return fmt.Errorf("Failed to download Mono FastCGI configuration: %v", err)
	}

	var config apps

	err = xml.Unmarshal([]byte(file), &config)

	if err != nil {
		return fmt.Errorf("Failed to parse Mono FastCGI configuration: %v", err)
	}

	for _, app := range config.Apps {
		if app.Name == cxt.domain {
			fmt.Println("Site found in Mono FastCGI configuration")
			return nil
		}
	}

	config.Apps = append(config.Apps, webApplication{cxt.domain, cxt.domain, 80, "/", "/var/www/" + cxt.domain, true})

	data, err := xml.MarshalIndent(&config, "", "	")

	if err != nil {
		return fmt.Errorf("Failed to generate xml for Mono FastCGI configuration: %v", err)
	}

	ssh.ScpUploadDataAsRoot(cxt.Client, string(data), "/etc/xsp4/debian.webapp", cxt.password)

	if err != nil {
		return fmt.Errorf("Failed to upload Mono FastCGI configuration: %v", err)
	}

	fmt.Println("Site added to Mono FastCGI configuration")

	return nil
}

type apps struct {
	Apps []webApplication `xml:"web-application"`
}

type webApplication struct {
	Name    string `xml:"name"`
	VHost   string `xml:"vhost"`
	VPort   int    `xml:"vport"`
	VPath   string `xml:"vpath"`
	Path    string `xml:"path"`
	Enabled bool   `xml:"enabled"`
}
