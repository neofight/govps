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

	pType, err := identifyProjectType()

	if err != nil {
		log.Fatal("Failed to identify a supported project type for deployment: ", err)
	}

	pipeline, err := createPipeline(pType)

	if err != nil {
		log.Fatal("Failed to create tasks for deployment: ", err)
	}

	password, err := promptForPassword()

	if err != nil {
		log.Fatal("Unable to read password: ", err)
	}

	client, err := ssh.CreateSSHClient(host, password)

	if err != nil {
		log.Fatal("Error creating SSH client: ", err)
	}

	defer client.Close()

	vps := Server{client: client, password: password}

	err = tasks.ExecutePipeline(tasks.Context{VPS: vps, Password: password, Domain: domain}, pipeline)

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
