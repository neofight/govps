package main

import (
	"fmt"
	"log"
	"os"

	"github.com/neofight/govps/ssh"
	"github.com/neofight/govps/tasks"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	host := os.Args[1]
	domain := os.Args[2]

	password, err := promptForPassword()

	if err != nil {
		log.Fatal("Unable to read password: ", err)
	}

	client, err := ssh.CreateSSHClient(host, password)

	if err != nil {
		log.Fatal("Error creating SSH client: ", err)
	}

	defer client.Close()

	pipeline := []tasks.Task{
		tasks.CreateMySQLDatabase{},
		tasks.PublishMVC{},
		tasks.AddSiteToMonoFastCGIConfiguration{},
		tasks.InstallMonoFastCGIService{},
		tasks.AddNginxFastCGIParameters{},
		tasks.AddNginxConfig{},
	}

	err = tasks.ExecutePipeline(tasks.Context{client, password, domain}, pipeline)

	if err != nil {
		log.Fatal("Error executing deployment process: ", err)
	}
}

func promptForPassword() ([]byte, error) {

	fmt.Print("Enter password:")
	password, err := terminal.ReadPassword(0)
	fmt.Println()

	return password, err
}
