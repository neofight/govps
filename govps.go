package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	host := os.Args[1]
	domain := os.Args[2]

	password, err := promptForPassword()

	if err != nil {
		log.Fatal("Unable to read password: ", err)
	}

	config, err := createNginxHostConfiguration(domain)

	if err != nil {
		log.Fatal("Error generating Nginx configuration file: ", err)
	}

	client, err := createSSHClient(host, password)

	if err != nil {
		log.Fatal("Error creating SSH client: ", err)
	}

	defer client.Close()

	err = scpUpload(client, config, domain)

	if err != nil {
		log.Fatal("Failed to upload file: ", err)
	}

	err = scpDownload(client, "/etc/nginx/fastcgi_params")

	if err != nil {
		log.Fatal("Failed to download file: ", err)
	}
}

func promptForPassword() ([]byte, error) {

	fmt.Print("Enter password:")
	password, err := terminal.ReadPassword(0)
	fmt.Println()

	return password, err
}
