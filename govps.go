package main

import (
	"log"
	"os"
)

func main() {
	host := os.Args[1]
	domain := os.Args[2]

	config, err := createNginxHostConfiguration(domain)

	if err != nil {
		log.Fatal("Error generating Nginx configuration file: ", err)
	}

	client, err := createSSHClient(host)

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
