package tasks

import (
	"bytes"
	"fmt"
	"text/template"
)

type Context struct {
	VPS      Server
	Password []byte
	Domain   string
}

type Task interface {
	Execute(cxt Context) error
}

func ExecutePipeline(cxt Context, tasks []Task) error {

	for _, task := range tasks {
		err := task.Execute(cxt)

		if err != nil {
			return err
		}
	}

	return nil
}

func uploadTemplate(cxt Context, templateName string, templateText string, templateData interface{}, remotePath string) error {

	template, err := template.New(templateName).Parse(templateText)

	if err != nil {
		return fmt.Errorf("Failed to parse template: %v", err)
	}

	var buffer bytes.Buffer

	template.Execute(&buffer, templateData)

	if err != nil {
		return fmt.Errorf("Failed to execute template: %v", err)
	}

	err = cxt.VPS.UploadData(buffer.String(), remotePath)

	if err != nil {
		return fmt.Errorf("Failed to upload file: %v", err)
	}

	return nil
}
