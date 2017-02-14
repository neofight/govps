package main

import (
	"os"
	"testing"

	"github.com/neofight/govps/io"
	"github.com/neofight/govps/mock"
)

func setupMockFile(file string) {

	io.FileSystem = mock.FileSystem{
		[]os.FileInfo{mock.FileInfo{FileName: file}},
	}
}

func TestIdentifyProjectTypeMvc(t *testing.T) {

	setupMockFile("Web.config")

	pType, err := identifyProjectType()

	if err != nil {
		t.Error("Expected project to be identfied but it was not")
	}

	if pType != Mvc {
		t.Error("Expected project to be identfied as MVC but it was not")
	}
}

func TestIdentifyProjectTypeStatic(t *testing.T) {

	setupMockFile("index.html")

	pType, err := identifyProjectType()

	if err != nil {
		t.Error("Expected project to be identfied but it was not")
	}

	if pType != Static {
		t.Error("Expected project to be identfied as MVC but it was not")
	}
}

func TestIdentifyProjectTypeUnknown(t *testing.T) {

	setupMockFile("xxx.xxx")

	pType, err := identifyProjectType()

	if err == nil {
		t.Error("Expected unknown project type to return an error but it did not")
	}

	if pType != Unknown {
		t.Error("Expected project to be identfied as Unknown but it was not")
	}
}

func TestIdentifyProjectTypeDirectoryError(t *testing.T) {

	io.FileSystem = mock.FileSystem{}

	pType, err := identifyProjectType()

	if err == nil {
		t.Error("Expected unknown project type to return an error but it did not")
	}

	if pType != Unknown {
		t.Error("Expected project to be identfied as Unknown but it was not")
	}
}

func TestCreatePipelineMvcErrorReadingPassword(t *testing.T) {

	io.Terminal = mock.Terminal{}

	pipeline, err := createPipeline(Mvc)

	if err == nil {
		t.Error("Expected failure to read password to return error but it did not")
	}

	if pipeline != nil {
		t.Error("Expected failure to read password to return a nil pipeline but it did not")
	}
}

func TestCreatePipelineMvcHappyPath(t *testing.T) {

	io.Terminal = mock.Terminal{"password"}

	pipeline, err := createPipeline(Mvc)

	if err != nil {
		t.Error("Expected pipeline to be created without error but it was not")
	}

	if pipeline == nil {
		t.Error("Expected pipeline to be created but it was not")
	}
}

func TestCreatePipelineStaticHappyPath(t *testing.T) {

	pipeline, err := createPipeline(Static)

	if err != nil {
		t.Error("Expected pipeline to be created without error but it was not")
	}

	if pipeline == nil {
		t.Error("Expected pipeline to be created but it was not")
	}
}

func TestCreatePipelineUnknown(t *testing.T) {

	pipeline, err := createPipeline(Unknown)

	if err == nil {
		t.Error("Expected unknown project type to return error it but did not")
	}

	if pipeline != nil {
		t.Error("Expected unknown project type to return a nil pipeline but it did not")
	}
}
