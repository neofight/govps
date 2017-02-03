package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/neofight/govps/tasks"
)

type projectType int

const (
	Unknown projectType = iota
	Mvc
	Static
)

func identifyProjectType() (projectType, error) {

	files, err := ioutil.ReadDir(".")

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

	return Unknown, fmt.Errorf("Failed to identify project type")
}

func createPipeline(pType projectType) ([]tasks.Task, error) {
	switch pType {
	case Mvc:
		return []tasks.Task{
			tasks.CreateMySQLDatabase{},
			tasks.PublishMVC{},
			tasks.AddSiteToMonoFastCGIConfiguration{},
			tasks.InstallMonoFastCGIService{},
			tasks.AddNginxFastCGIParameters{},
			tasks.AddNginxConfig{},
		}, nil
	case Static:
		return []tasks.Task{
			tasks.PublishStatic{},
		}, nil
	default:
		return nil, fmt.Errorf("Project type %v is not supported", pType)
	}
}
