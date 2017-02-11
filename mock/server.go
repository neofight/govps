package mock

import (
	"fmt"
	"os"
)

type Server struct {
	UploadedData string
	UploadedPath string

	SudoCommandsRun []string
}

func (s *Server) RunCommand(command string, inputs ...string) (string, error) {
	return "", nil
}

func (s *Server) RunSudoCommands(commands ...string) (string, error) {

	s.SudoCommandsRun = append(s.SudoCommandsRun, commands...)

	return "", nil
}

func (s *Server) DownloadFile(path string) (string, error) {
	return "", nil
}

func (s *Server) UploadFiles(localPath string, remotePath string, filter func(path string, info os.FileInfo) bool) error {
	return nil
}

func (s *Server) UploadData(data string, filePath string) error {

	fmt.Println(s.UploadedPath)

	s.UploadedData = data
	s.UploadedPath = filePath

	return nil
}
