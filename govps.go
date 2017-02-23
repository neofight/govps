package main

import (
	"log"
	"os"

	"github.com/neofight/govps/ssh"
	"github.com/neofight/govps/tasks"
)

func main() {

	args, ok := Parse(os.Args)

	if !ok {
		return;
	}

	pType, err := IdentifyProjectType()

	if err != nil {
		log.Fatal("Failed to identify a supported project type for deployment: ", err)
	}

	password, err := promptForPassword("server")

	if err != nil {
		log.Fatal("Unable to read password: ", err)
	}

	pipeline, err := CreatePipeline(pType)

	if err != nil {
		log.Fatal("Failed to create tasks for deployment: ", err)
	}

	client, err := ssh.CreateSSHClient(args.Host, password)

	if err != nil {
		log.Fatal("Error creating SSH client: ", err)
	}

	defer client.Close()

	vps := Server{client: client, password: password}

	err = tasks.ExecutePipeline(tasks.Context{VPS: vps, Domain: args.Domain}, pipeline)

	if err != nil {
		log.Fatal("Error executing deployment process: ", err)
	}
}
