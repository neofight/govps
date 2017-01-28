package main

import (
	"bytes"
	"fmt"
	"text/template"

	ssh "github.com/neofight/govps/ssh"
	cryptossh "golang.org/x/crypto/ssh"
)

type context struct {
	Client   *cryptossh.Client
	password []byte
	domain   string
}

type step interface {
	Execute(cxt context) error
}

func executePipeline(cxt context, steps []step) error {

	for _, step := range steps {
		err := step.Execute(cxt)

		if err != nil {
			return err
		}
	}

	return nil
}

func uploadTemplate(cxt context, templateName string, templateText string, templateData interface{}, remotePath string) error {

	template, err := template.New(templateName).Parse(templateText)

	if err != nil {
		return fmt.Errorf("Failed to parse template: %v", err)
	}

	var buffer bytes.Buffer

	template.Execute(&buffer, templateData)

	if err != nil {
		return fmt.Errorf("Failed to execute template: %v", err)
	}

	err = ssh.ScpUploadDataAsRoot(cxt.Client, buffer.String(), remotePath, cxt.password)

	if err != nil {
		return fmt.Errorf("Failed to upload file: %v", err)
	}

	return nil
}
