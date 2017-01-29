package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neofight/govps/ssh"
)

type publishMVC struct {
}

func (step publishMVC) Execute(cxt context) error {

	err := ssh.RunSudoCommand(cxt.Client, "mkdir -p /var/www/"+cxt.domain, cxt.password)

	if err != nil {
		return fmt.Errorf("Unable to create directory: %v", err)
	}

	return ssh.ScpUploadAsRoot(cxt.Client, ".", "/var/www/"+cxt.domain, cxt.password, func(path string, info os.FileInfo) bool {

		if info.IsDir() {
			switch path {
			case "obj", "packages":
				return false
			default:
				return true
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
