package main

import "fmt"
import "log"
import "os"

func main() {
	host := os.Args[1]

	client, err := createSSHClient(host)

	if err != nil {
		log.Fatal("Error creating SSH client: ", err)
	}

	defer client.close()

	result, err := client.runCommands("uname -a")

	if err != nil {
		log.Fatal("Failed to run remote command: ", err)
	}

	fmt.Print(result)
}
