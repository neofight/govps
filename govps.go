package main

import (
	"fmt"
	"log"
	"os"

	"github.com/neofight/govps/ssh"
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

	pipeline := []step{
		addNginxConfig{},
		addNginxFastCGIParameters{},
		publishMVC{},
		installMonoFastCGIService{},
		createMySQLDatabase{},
	}

	err = executePipeline(context{client, password, domain}, pipeline)

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
