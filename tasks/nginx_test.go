package tasks

import (
	"strings"
	"testing"

	"github.com/neofight/govps/mock"
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

	server.Files[fastcgiParamsPath] = fastcgiParamsFile

	cxt := Context{server, "test.com"}

	var task = AddNginxFastCGIParameters{}

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected Fast CGI Parameters to be added without error but they were not")
	}

	if server.UploadedPath != fastcgiParamsPath {
		t.Errorf("Expected the remote path to be as follows:\n%v\nBut was:\n%v", fastcgiParamsPath, server.UploadedPath)
	}

	if c := strings.Count(server.UploadedData, pathInfoParameter); c != 1 {
		t.Errorf("Expected the uploaded configuration file to include the PATH INFO setting once but was found %v times", c)
	}

	if c := strings.Count(server.UploadedData, scriptFilenameParameter); c != 1 {
		t.Errorf("Expected the uploaded configuration file to include the SCRIPT FILENAME setting once but was found %v times", c)
	}
}

func TestExistingNginxFastCGIParameters(t *testing.T) {

	server := mock.NewServer()

	server.Files[fastcgiParamsPath] = fastcgiParamsFile + pathInfoParameter + scriptFilenameParameter

	cxt := Context{server, "test.com"}

	var task = AddNginxFastCGIParameters{}

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected Fast CGI Parameters to be added without error but they were not")
	}

	if server.UploadedPath != "" || server.UploadedData != "" {
		t.Errorf("Expected the configuration file not to be uploaded but it was")
	}
}
