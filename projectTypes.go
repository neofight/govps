package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/neofight/govps/io"
	"github.com/neofight/govps/tasks"
)

type projectType int

const (
	Unknown projectType = iota
	Mvc
	Static
)

func IdentifyProjectType() (projectType, error) {

	files, err := io.FileSystem.ReadDir(".")

	if err != nil {
		return Unknown, fmt.Errorf("Failed to read current directory: %v", err)
	}

	for _, file := range files {

		switch {
		case file.Name() == "Web.config":
			return Mvc, nil
		case filepath.Ext(file.Name()) == ".html":
			return Static, nil
		}
	}

	return Unknown, errors.New("Failed to identify project type")
}

func CreatePipeline(pType projectType) ([]tasks.Task, error) {
	switch pType {
	case Mvc:
		mySQLPassword, err := promptForPassword("MySql")

		if err != nil {
			return nil, fmt.Errorf("Failed to read password: %v", err)
		}

		return []tasks.Task{
			tasks.CreateMySQLDatabase{MySQLPassword: mySQLPassword},
			tasks.PublishMVC{},
			tasks.ConfigureMonoSite{},
			tasks.InstallMono{},
			tasks.AddNginxFastCGIParameters{},
			tasks.UploadMvcNginxConfig{},
			tasks.EnableNginxSite{},
		}, nil
	case Static:
		return []tasks.Task{
			tasks.PublishStatic{},
			tasks.UploadStaticNginxConfig{},
			tasks.EnableNginxSite{},
		}, nil
	default:
		return nil, fmt.Errorf("Project type %v is not supported", pType)
	}
}
