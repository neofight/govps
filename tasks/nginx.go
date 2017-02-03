package tasks

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/neofight/govps/ssh"
)

type AddNginxConfig struct {
}

func (step AddNginxConfig) Execute(cxt Context) error {

	err := uploadTemplate(cxt, "nginx", nginxTemplate, cxt.Domain, "/etc/nginx/sites-available/"+cxt.Domain)

	if err != nil {
		return fmt.Errorf("Failed to deploy nginx configuration file for %v: %v", cxt.Domain, err)
	}

	fmt.Printf("Nginx configuration file for %v uploaded\n", cxt.Domain)

	return nil
}

type AddNginxFastCGIParameters struct {
}

func (step AddNginxFastCGIParameters) Execute(cxt Context) error {

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

	successMessage := ""

	if !pathInfoSet {
		file += "\nfastcgi_param  PATH_INFO          \"\";"
		successMessage += "Added PATH_INFO to fastcgi_params\n"
	} else {
		successMessage += "PATH_INFO already present in fastcgi_params\n"
	}

	if !scriptFilenameSet {
		file += "\nfastcgi_param  SCRIPT_FILENAME    $document_root$fastcgi_script_name;"
		successMessage += "Added SCRIPT_FILENAME to fastcgi_params\n"
	} else {
		successMessage += "SCRIPT_FILENAME already present in fastcgi_params\n"
	}

	if !pathInfoSet || !scriptFilenameSet {

		err = ssh.ScpUploadDataAsRoot(cxt.Client, file, "/etc/nginx/fastcgi_params", cxt.Password)

		if err != nil {
			return fmt.Errorf("Failed to update fastcgi_params: %v", err)
		}
	}

	fmt.Print(successMessage)

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