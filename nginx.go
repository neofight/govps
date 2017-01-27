package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/neofight/govps/ssh"
)

type addNginxConfig struct {
}

func (step addNginxConfig) Execute(cxt context) error {

	template, err := template.New("nginx").Parse(nginxTemplate)

	if err != nil {
		return fmt.Errorf("Unable to parse Nginx Host Configuration template: %v", err)
	}

	var buffer bytes.Buffer

	template.Execute(&buffer, cxt.domain)

	err = ssh.ScpUploadDataAsRoot(cxt.Client, buffer.String(), "/etc/nginx/sites-available/"+cxt.domain, cxt.password)

	if err != nil {
		return fmt.Errorf("Failed to upload file: %v", err)
	}

	fmt.Printf("Nginx configuration file for %v uploaded\n", cxt.domain)

	return nil
}

type addNginxFastCGIParameters struct {
}

func (step addNginxFastCGIParameters) Execute(cxt context) error {

	file, err := ssh.ScpDownloadFile(cxt.Client, "/etc/nginx/fastcgi_params")

	if err != nil {
		return fmt.Errorf("Failed to download file: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(file))

	pathInfoSet := false
	scriptFilenameSet := false

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())

		if len(parts) == 3 && parts[0] == "fastcgi_param" {
			switch parts[1] {
			case "PATH_INFO":
				pathInfoSet = true
			case "SCRIPT_FILENAME":
				scriptFilenameSet = true
			}
		}
	}

	if !pathInfoSet {
		file += "\nfastcgi_param  PATH_INFO          \"\";"
		fmt.Println("Added PATH_INFO to fastcgi_params")
	} else {
		fmt.Println("PATH_INFO already present in fastcgi_params")
	}

	if !scriptFilenameSet {
		file += "\nfastcgi_param  SCRIPT_FILENAME    $document_root$fastcgi_script_name;"
		fmt.Println("Added SCRIPT_FILENAME to fastcgi_params")
	} else {
		fmt.Println("SCRIPT_FILENAME already present in fastcgi_params")
	}

	if pathInfoSet && scriptFilenameSet {
		return nil
	} else {
		return ssh.ScpUploadDataAsRoot(cxt.Client, file, "/etc/nginx/fastcgi_params", cxt.password)
	}
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
