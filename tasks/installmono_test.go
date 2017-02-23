package tasks_test

import (
	"testing"

	"github.com/neofight/govps/mock"
	"github.com/neofight/govps/tasks"
)

func executeInstallMonoFastCGIService(server tasks.Server) error {

	cxt := tasks.Context{server, "test.com"}

	var task tasks.InstallMono

	return task.Execute(cxt)
}

func TestInstallMonoFastCGIServiceUnitFile(t *testing.T) {

	server := mock.NewServer()

	err := executeInstallMonoFastCGIService(server)

	if err != nil {
		t.Error("Expected unit file to be uploaded without error but it was not")
	}

	if server.UploadedData != tasks.MonoUnitFile {
		t.Errorf("Expected the uploaded unit file to be as follows:\n%v\nBut was:\n%v", tasks.MonoUnitFile, server.UploadedData)
	}

	if server.UploadedPath != tasks.MonoUnitFilePath {
		t.Errorf("Expected the remote path to be as follows:\n%v\nBut was:\n%v", tasks.MonoUnitFilePath, server.UploadedPath)
	}
}

func TestInstallMonoFastCGIServiceEnableService(t *testing.T) {

	server := mock.NewServer()

	err := executeInstallMonoFastCGIService(server)

	if err != nil {
		t.Error("Expected service to be enabled without error but it was not")
	}

	if !contains(server.SudoCommandsRun, tasks.EnableMonoServiceCommand) {
		t.Error("Expected service to be enabled but it was not")
	}
}

func TestInstallMonoFastCGIServiceStartService(t *testing.T) {

	server := mock.NewServer()

	err := executeInstallMonoFastCGIService(server)

	if err != nil {
		t.Error("Expected service to be enabled without error but it was not")
	}

	if !contains(server.SudoCommandsRun, tasks.StartMonoServiceCommand) {
		t.Error("Expected service to be started but it was not")
	}
}
