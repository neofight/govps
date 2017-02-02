package tasks

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/neofight/govps/ssh"
)

type CreateMySQLDatabase struct {
}

type connectionStrings struct {
	ConnectionStrings []add `xml:"connectionStrings>add"`
}

type add struct {
	Name             string `xml:"name,attr"`
	ConnectionString string `xml:"connectionString,attr"`
	ProviderName     string `xml:"providerName,attr"`
}

type database struct {
	Name     string
	User     string
	Password string
}

func (step CreateMySQLDatabase) Execute(cxt Context) error {

	var db database

	file, err := os.Open("Web.config")

	if err != nil {
		return fmt.Errorf("Failed to open Web.config: %v", err)
	}

	decoder := xml.NewDecoder(file)

	for {
		token, err := decoder.Token()

		if err != nil {
			break
		}

		switch tag := token.(type) {
		case xml.StartElement:
			if tag.Name.Local != "connectionStrings" {
				var c connectionStrings

				decoder.DecodeElement(&c, &tag)

				pairs := strings.Split(c.ConnectionStrings[0].ConnectionString, ";")

				for _, pair := range pairs {
					parts := strings.Split(pair, "=")

					switch parts[0] {
					case "Database":
						db.Name = parts[1]
					case "Uid":
						db.User = parts[1]
					case "Pwd":
						db.Password = parts[1]
					}
				}

				break
			}
		}
	}

	checkCommand := fmt.Sprintf("mysql -u root -p -BNe \"SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%v'\"", db.Name)
	password := fmt.Sprintf("%v\n", string(cxt.Password))

	result, err := ssh.RunCommand(cxt.Client, cxt.Password, checkCommand, password)

	if err != nil {
		return fmt.Errorf("Failed to run query: %v", err)
	}

	if strings.TrimSpace(result) == db.Name {
		fmt.Println("Database already created")
		return nil
	}

	create := fmt.Sprintf("CREATE DATABASE %v;\n", db.Name)
	grant := fmt.Sprintf("GRANT ALL ON %v.* TO '%v' IDENTIFIED BY '%v';\n", db.Name, db.User, db.Password)
	quit := "QUIT\n"

	_, err = ssh.RunCommand(cxt.Client, cxt.Password, "mysql -u root -p", password, create, grant, quit)

	if err != nil {
		return fmt.Errorf("Failed to create database %v: %v", db.Name, err)
	}

	fmt.Println("Database created")

	return nil
}
