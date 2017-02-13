package tasks

import (
	"github.com/neofight/govps/mock"
	"testing"
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

	server.Files[monoSiteConfigurationPath] = emptyMonoSiteConfiguration

	cxt := Context{server, "test.com"}

	var task ConfigureMonoSite

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected configuration file to be uploaded without error but it was not")
	}

	if server.UploadedData != testMonoSiteConfiguration {
		t.Errorf("Expected the uploaded configuration file to be as follows:\n%v\nBut was:\n%v", testMonoSiteConfiguration, server.UploadedData)
	}

	if server.UploadedPath != monoSiteConfigurationPath {
		t.Errorf("Expected the remote path to be as follows:\n%v\nBut was:\n%v", monoSiteConfigurationPath, server.UploadedPath)
	}
}

func TestConfigureMonoSiteExistingConfiguration(t *testing.T) {

	server := mock.NewServer()

	server.Files[monoSiteConfigurationPath] = testMonoSiteConfiguration

	cxt := Context{server, "test.com"}

	var task ConfigureMonoSite

	err := task.Execute(cxt)

	if err != nil {
		t.Error("Expected the task to complete without error but it did not")
	}

	if server.UploadedPath != "" || server.UploadedData != "" {
		t.Errorf("Expected the configuration file not to be uploaded but it was")
	}
}
