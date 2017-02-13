package mock

import (
	"errors"
	"fmt"
	"os"
)

type Server struct {
	Files map[string]string

	UploadedData string
	UploadedPath string

	SudoCommandsRun []string
}

func NewServer() *Server {
	return &Server{Files: make(map[string]string)}
}

func (s *Server) RunCommand(command string, inputs ...string) (string, error) {
	return "", nil
}

func (s *Server) RunSudoCommands(commands ...string) (string, error) {

	s.SudoCommandsRun = append(s.SudoCommandsRun, commands...)

	return "", nil
}

func (s *Server) DownloadFile(path string) (string, error) {

	data, ok := s.Files[path]

	if !ok {
		return "", errors.New("File does not exist")
	}

	return data, nil
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
