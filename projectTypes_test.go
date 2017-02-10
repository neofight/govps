package main

import (
	"os"
	"testing"

	"github.com/neofight/govps/mock"
)

func setupMockFile(file string) {

	fs = mock.FileSystem{
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

	fs = mock.FileSystem{}

	pType, err := identifyProjectType()

	if err == nil {
		t.Error("Expected unknown project type to return an error but it did not")
	}

	if pType != Unknown {
		t.Error("Expected project to be identfied as Unknown but it was not")
	}
}
