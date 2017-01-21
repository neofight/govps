package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type AddNginxConfig struct {
	domain   string
	password []byte
}

func (step AddNginxConfig) Execute(cxt context) error {

	template, err := template.New("nginx").Parse(nginxTemplate)

	if err != nil {
		return fmt.Errorf("Unable to parse Nginx Host Configuration template: %v", err)
	}

	var buffer bytes.Buffer

	template.Execute(&buffer, step.domain)

	err = scpUpload(cxt.Client, buffer.String(), step.domain)

	if err != nil {
		return fmt.Errorf("Failed to upload file: %v", err)
	}

	session, err := cxt.Client.NewSession()

	if err != nil {
		return fmt.Errorf("Unable to create session: %v", err)
	}

	defer session.Close()

	_, err = runSudoCommands(session, step.password, "mv "+step.domain+" /etc/nginx/sites-available/")

	if err != nil {
		return fmt.Errorf("Unable to move Nginx configuration file to the correct location: %v", err)
	}

	return nil
}

type AddNginxFastCGIParameters struct {
}

func (step AddNginxFastCGIParameters) Execute(cxt context) error {

	data, err := scpDownload(cxt.Client, "/etc/nginx/fastcgi_params")

	if err != nil {
		return fmt.Errorf("Failed to download file: %v", err)
	}

	fmt.Print(data)

	return nil
}

var nginxTemplate = `server {
        listen   80;
        server_name  {{.}};
        access_log   /var/log/nginx/{{.}}.access.log;
        root /var/www/{{.}}/;

        location / {
                index index.html index.htm default.aspx Default.aspx;
                fastcgi_index Home;
                fastcgi_pass 127.0.0.1:9000;
                include /etc/nginx/fastcgi_params;
        }
}`
