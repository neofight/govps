package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neofight/govps/ssh"
)

type PublishMVC struct {
}

func (step PublishMVC) Execute(cxt Context) error {

	deployPath := "/var/www/" + cxt.Domain

	_, err := ssh.RunSudoCommands(cxt.Client, cxt.Password, "rm -rf "+deployPath, "mkdir "+deployPath)

	if err != nil {
		return fmt.Errorf("Unable to create directory: %v", err)
	}

	err = ssh.ScpUploadAsRoot(cxt.Client, ".", deployPath, cxt.Password, func(path string, info os.FileInfo) bool {

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

	if err != nil {
		return fmt.Errorf("Failed to upload website: %v", err)
	}

	fmt.Println("Website published")

	return nil
}
