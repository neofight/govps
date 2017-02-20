package tasks

import (
	"bufio"
	"fmt"
	"strings"
)

type AddNginxFastCGIParameters struct {
}

func (task AddNginxFastCGIParameters) Execute(cxt Context) error {

	file, err := cxt.VPS.DownloadFile(fastcgiParamsPath)

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
		file += pathInfoParameter
		successMessage += "Added PATH_INFO to fastcgi_params\n"
	} else {
		successMessage += "PATH_INFO already present in fastcgi_params\n"
	}

	if !scriptFilenameSet {
		file += scriptFilenameParameter
		successMessage += "Added SCRIPT_FILENAME to fastcgi_params\n"
	} else {
		successMessage += "SCRIPT_FILENAME already present in fastcgi_params\n"
	}

	if !pathInfoSet || !scriptFilenameSet {

		err = cxt.VPS.UploadData(file, "/etc/nginx/fastcgi_params")

		if err != nil {
			return fmt.Errorf("Failed to update fastcgi_params: %v", err)
		}
	}

	fmt.Print(successMessage)

	return nil
}

const fastcgiParamsPath = "/etc/nginx/fastcgi_params"

const pathInfoParameter = "\nfastcgi_param  PATH_INFO          \"\";"
const scriptFilenameParameter = "\nfastcgi_param  SCRIPT_FILENAME    $document_root$fastcgi_script_name;"

func uploadNginxConfig(cxt Context, name string, templateText string) error {

	err := uploadTemplate(cxt, name, templateText, cxt.Domain, sitesAvailablePath+cxt.Domain)

	if err != nil {
		return fmt.Errorf("Failed to deploy nginx configuration file for %v: %v", cxt.Domain, err)
	}

	fmt.Printf("Nginx configuration file for %v uploaded\n", cxt.Domain)

	return nil
}

type UploadMvcNginxConfig struct {
}

func (task UploadMvcNginxConfig) Execute(cxt Context) error {

	return uploadNginxConfig(cxt, "mvcNginx", mvcTemplate)
}

var mvcTemplate = `server {
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

type UploadStaticNginxConfig struct {
}

func (task UploadStaticNginxConfig) Execute(cxt Context) error {

	return uploadNginxConfig(cxt, "staticNginx", staticTemplate)
}

var staticTemplate = `server {
        listen		80;
        server_name	{{.}};
        access_log	/var/log/nginx/{{.}}.access.log;
        root		/var/www/{{.}}/;

        location / {
                index			index.html;
                error_page	404	/404.html;
        }
}`

type EnableNginxSite struct {
}

func (task EnableNginxSite) Execute(cxt Context) error {

	enableSite := fmt.Sprintf(enableSiteCommand, cxt.Domain, cxt.Domain)

	reloadConfig := reloadConfigCommand

	_, err := cxt.VPS.RunSudoCommands(enableSite, reloadConfig)

	if err != nil {
		return fmt.Errorf("Failed to enable site: %v", err)
	}

	fmt.Println("Site enabled")

	return nil
}

const sitesAvailablePath = "/etc/nginx/sites-available/"
const enableSiteCommand = "ln -sf /etc/nginx/sites-available/%v /etc/nginx/sites-enabled/%v"
const reloadConfigCommand = "systemctl reload nginx"
