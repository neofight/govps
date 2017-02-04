package tasks

import (
	"os"
)

type Server interface {
	RunCommand(command string, inputs ...string) (string, error)
	RunSudoCommands(commands ...string) (string, error)
	DownloadFile(path string) (string, error)
	UploadFiles(localPath string, remotePath string, filter func(path string, info os.FileInfo) bool) error
	UploadData(data string, filePath string) error
}
