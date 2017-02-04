package main

import (
	"os"

	"github.com/neofight/govps/ssh"
	cryptossh "golang.org/x/crypto/ssh"
)

type Server struct {
	client   *cryptossh.Client
	password []byte
}

func (server Server) RunCommand(command string, inputs ...string) (string, error) {

	return ssh.RunCommand(server.client, server.password, command, inputs...)
}

func (server Server) RunSudoCommands(commands ...string) (string, error) {

	return ssh.RunSudoCommands(server.client, server.password, commands...)
}

func (server Server) DownloadFile(path string) (string, error) {

	return ssh.ScpDownloadFile(server.client, path)
}

func (server Server) UploadFiles(localPath string, remotePath string, filter func(path string, info os.FileInfo) bool) error {

	return ssh.ScpUpload(server.client, localPath, remotePath, server.password, filter)
}

func (server Server) UploadData(data string, filePath string) error {

	return ssh.ScpUploadData(server.client, data, filePath, server.password)
}
