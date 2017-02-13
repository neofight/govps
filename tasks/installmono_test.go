package tasks

import (
	"testing"

	"github.com/neofight/govps/mock"
)

func contains(a []string, v string) bool {

	for _, av := range a {
		if av == v {
			return true
		}
	}

	return false
}

func executeInstallMonoFastCGIService(server Server) error {

	cxt := Context{server, "test.com"}

	var task InstallMono

	return task.Execute(cxt)
}

func TestInstallMonoFastCGIServiceUnitFile(t *testing.T) {

	server := mock.NewServer()

	err := executeInstallMonoFastCGIService(server)

	if err != nil {
		t.Error("Expected unit file to be uploaded without error but it was not")
	}

	if server.UploadedData != monoUnitFile {
		t.Errorf("Expected the uploaded unit file to be as follows:\n%v\nBut was:\n%v", monoUnitFile, server.UploadedData)
	}

	if server.UploadedPath != monoUnitFilePath {
		t.Errorf("Expected the remote path to be as follows:\n%v\nBut was:\n%v", monoUnitFilePath, server.UploadedPath)
	}
}

func TestInstallMonoFastCGIServiceEnableService(t *testing.T) {

	server := mock.NewServer()

	err := executeInstallMonoFastCGIService(server)

	if err != nil {
		t.Error("Expected service to be enabled without error but it was not")
	}

	if !contains(server.SudoCommandsRun, enableMonoServiceCommand) {
		t.Error("Expected service to be enabled but it was not")
	}
}

func TestInstallMonoFastCGIServiceStartService(t *testing.T) {

	server := mock.NewServer()

	err := executeInstallMonoFastCGIService(server)

	if err != nil {
		t.Error("Expected service to be enabled without error but it was not")
	}

	if !contains(server.SudoCommandsRun, startMonoServiceCommand) {
		t.Error("Expected service to be started but it was not")
	}
}
