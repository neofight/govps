package tasks_test

import (
	"testing"

	"github.com/neofight/govps/mock"
	"github.com/neofight/govps/tasks"
)

const emptyMonoSiteConfiguration = "<apps></apps>"

const testMonoSiteConfiguration = `<apps>
	<web-application>
		<name>test.com</name>
		<vhost>test.com</vhost>
		<vport>80</vport>
		<vpath>/</vpath>
		<path>/var/www/test.com</path>
		<enabled>true</enabled>
	</web-application>
</apps>`

func TestConfigureMonoSiteNoConfiguration(t *testing.T) {

	server := mock.NewServer()

	server.Files[tasks.MonoSiteConfigurationPath] = emptyMonoSiteConfiguration

	cxt := tasks.Context{VPS: server, Domain: "test.com"}

	var task tasks.ConfigureMonoSite

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected configuration file to be uploaded without error but it was not")
	}

	uploadedData, ok := server.UploadedFiles[tasks.MonoSiteConfigurationPath]

	if !ok {
		t.Errorf("Expected the configuration to be uploaded to:\n%v", tasks.MonoSiteConfigurationPath)
		return
	}

	if uploadedData != testMonoSiteConfiguration {
		t.Errorf("Expected the uploaded configuration file to be as follows:\n%v\nBut was:\n%v", testMonoSiteConfiguration, uploadedData)
	}
}

func TestConfigureMonoSiteExistingConfiguration(t *testing.T) {

	server := mock.NewServer()

	server.Files[tasks.MonoSiteConfigurationPath] = testMonoSiteConfiguration

	cxt := tasks.Context{VPS: server, Domain: "test.com"}

	var task tasks.ConfigureMonoSite

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected the task to complete without error but it did not")
	}

	_, ok := server.UploadedFiles[tasks.MonoSiteConfigurationPath]

	if ok {
		t.Error("Expected the configuration file not to be uploaded but it was")
	}
}
