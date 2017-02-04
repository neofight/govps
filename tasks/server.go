package tasks

import (
	"os"
)

type Server interface {
	RunCommand(command string, inputs ...string) (string, error)
	RunSudoCommands(commands ...string) (string, error)
	ScpDownloadFile(path string) (string, error)
	ScpUploadAsRoot(localPath string, remotePath string, filter func(path string, info os.FileInfo) bool) error
	ScpUploadDataAsRoot(data string, filePath string) error
}
