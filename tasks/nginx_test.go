package tasks_test

import (
	"strings"
	"testing"

	"fmt"
	"github.com/neofight/govps/mock"
	"github.com/neofight/govps/tasks"
)

const fastcgiParamsFile = `fastcgi_param  QUERY_STRING       $query_string;
fastcgi_param  REQUEST_METHOD     $request_method;
fastcgi_param  CONTENT_TYPE       $content_type;
fastcgi_param  CONTENT_LENGTH     $content_length;

fastcgi_param  SCRIPT_NAME        $fastcgi_script_name;
fastcgi_param  REQUEST_URI        $request_uri;
fastcgi_param  DOCUMENT_URI       $document_uri;
fastcgi_param  DOCUMENT_ROOT      $document_root;
fastcgi_param  SERVER_PROTOCOL    $server_protocol;
fastcgi_param  REQUEST_SCHEME     $scheme;
fastcgi_param  HTTPS              $https if_not_empty;

fastcgi_param  GATEWAY_INTERFACE  CGI/1.1;
fastcgi_param  SERVER_SOFTWARE    nginx/$nginx_version;

fastcgi_param  REMOTE_ADDR        $remote_addr;
fastcgi_param  REMOTE_PORT        $remote_port;
fastcgi_param  SERVER_ADDR        $server_addr;
fastcgi_param  SERVER_PORT        $server_port;
fastcgi_param  SERVER_NAME        $server_name;

# PHP only, required if PHP was built with --enable-force-cgi-redirect
fastcgi_param  REDIRECT_STATUS    200;`

func TestAddNginxFastCGIParameters(t *testing.T) {

	server := mock.NewServer()

	server.Files[tasks.FastcgiParamsPath] = fastcgiParamsFile

	cxt := tasks.Context{VPS: server, Domain: "test.com"}

	var task = tasks.AddNginxFastCGIParameters{}

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected Fast CGI Parameters to be added without error but they were not")
	}

	uploadedData, ok := server.UploadedFiles[tasks.FastcgiParamsPath]

	if !ok {
		t.Errorf("Expected the configuration to be uploaded to:\n%v", tasks.FastcgiParamsPath)
		return
	}

	if c := strings.Count(uploadedData, tasks.PathInfoParameter); c != 1 {
		t.Errorf("Expected the uploaded configuration file to include the PATH INFO setting once but was found %v times", c)
	}

	if c := strings.Count(uploadedData, tasks.ScriptFilenameParameter); c != 1 {
		t.Errorf("Expected the uploaded configuration file to include the SCRIPT FILENAME setting once but was found %v times", c)
	}
}

func TestExistingNginxFastCGIParameters(t *testing.T) {

	server := mock.NewServer()

	server.Files[tasks.FastcgiParamsPath] = fastcgiParamsFile + tasks.PathInfoParameter + tasks.ScriptFilenameParameter

	cxt := tasks.Context{VPS: server, Domain: "test.com"}

	var task = tasks.AddNginxFastCGIParameters{}

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected Fast CGI Parameters to be added without error but they were not")
	}

	_, ok := server.UploadedFiles[tasks.FastcgiParamsPath]

	if ok {
		t.Error("Expected the configuration file not to be uploaded but it was")
	}
}

const testMvcNginxConfiguration = `server {
        listen   80;
        server_name  test.com;
        access_log   /var/log/nginx/test.com.access.log;
        root /var/www/test.com/;

        location / {
                index index.html index.htm default.aspx Default.aspx;
                fastcgi_index Home;
                fastcgi_pass 127.0.0.1:9000;
                include /etc/nginx/fastcgi_params;
        }
}`

func TestUploadMvcNginxConfig(t *testing.T) {

	server := mock.NewServer()

	cxt := tasks.Context{VPS: server, Domain: "test.com"}

	var task = tasks.UploadMvcNginxConfig{}

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected MVC Nginx Config to uploaded without error but it was not")
	}

	uploadedData, ok := server.UploadedFiles[tasks.SitesAvailablePath+"test.com"]

	if !ok {
		t.Errorf("Expected the configuration to be uploaded to:\n%v", tasks.SitesAvailablePath+"test.com")
	}

	if uploadedData != testMvcNginxConfiguration {
		t.Errorf("Expected the uploaded configuration file to be as follows:\n%v\nBut was:\n%v", testMvcNginxConfiguration, uploadedData)
	}
}

const testStaticNginxConfiguration = `server {
        listen		80;
        server_name	test.com;
        access_log	/var/log/nginx/test.com.access.log;
        root		/var/www/test.com/;

        location / {
                index			index.html;
                error_page	404	/404.html;
        }
}`

func TestUploadStaticNginxConfig(t *testing.T) {

	server := mock.NewServer()

	cxt := tasks.Context{VPS: server, Domain: "test.com"}

	var task = tasks.UploadStaticNginxConfig{}

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected Static Nginx Config to uploaded without error but it was not")
	}

	uploadedData, ok := server.UploadedFiles[tasks.SitesAvailablePath+"test.com"]

	if !ok {
		t.Errorf("Expected the configuration to be uploaded to:\n%v", tasks.SitesAvailablePath+"test.com")
	}

	if uploadedData != testStaticNginxConfiguration {
		t.Errorf("Expected the uploaded configuration file to be as follows:\n%v\nBut was:\n%v", testStaticNginxConfiguration, uploadedData)
	}
}

func positionOfCommand(server *mock.Server, command string) (index int, ok bool) {

	for i, runCommand := range server.SudoCommandsRun {

		if runCommand == command {
			return i, true
		}
	}

	return 0, false
}

func TestEnableNginxSite(t *testing.T) {

	server := mock.NewServer()

	cxt := tasks.Context{VPS: server, Domain: "test.com"}

	var task = tasks.EnableNginxSite{}

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected Nginx site to be enabled without error but it was not")
	}

	enablePosition, ok := positionOfCommand(server, fmt.Sprintf(tasks.EnableSiteCommand, "test.com", "test.com"))

	if !ok {
		t.Error("Expected site to be enabled but it was not")
	}

	reloadPosition, ok := positionOfCommand(server, tasks.ReloadConfigCommand)

	if !ok {
		t.Error("Expected Nginx configuration to be reloaded but it was not")
	}

	if enablePosition > reloadPosition {
		t.Error("The Nginx configuration was reloaded before the site was enabled")
	}
}
