package tasks

import (
	"fmt"
	"os"
	"path/filepath"
)

type PublishMVC struct {
}

func (task PublishMVC) Execute(cxt Context) error {

	return publish(cxt, func(path string, info os.FileInfo) bool {
		if info.IsDir() {
			if path != filepath.Base(path) {
				return true
			}

			switch path {
			case ".", "bin", "Content", "fonts", "Scripts", "Views":
				return true
			default:
				return false
			}
		}

		dir, _ := filepath.Split(path)

		if dir == "" {
			switch path {
			case "Global.asax", "Web.config":
				return true
			default:
				return false
			}
		}

		return true
	})
}

type PublishStatic struct {
}

func (task PublishStatic) Execute(cxt Context) error {

	return publish(cxt, func(path string, info os.FileInfo) bool {
		return true
	})
}

func publish(cxt Context, filter func(path string, info os.FileInfo) bool) error {

	deployPath := "/var/www/" + cxt.Domain

	_, err := cxt.VPS.RunSudoCommands("rm -rf "+deployPath, "mkdir "+deployPath)

	if err != nil {
		return fmt.Errorf("Unable to create directory: %v", err)
	}

	err = cxt.VPS.UploadFiles(".", deployPath, filter)

	if err != nil {
		return fmt.Errorf("Failed to upload website: %v", err)
	}

	fmt.Println("Website published")

	return nil
}
