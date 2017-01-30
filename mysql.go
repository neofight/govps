package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/neofight/govps/ssh"
)

type createMySQLDatabase struct {
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

func (step createMySQLDatabase) Execute(cxt context) error {

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

	password := fmt.Sprintf("%v\n", string(cxt.password))
	create := fmt.Sprintf("CREATE DATABASE %v;\n", db.Name)
	grant := fmt.Sprintf("GRANT ALL ON %v.* TO '%v' IDENTIFIED BY '%v';\n", db.Name, db.User, db.Password)
	quit := "QUIT\n"

	return ssh.RunCommand(cxt.Client, cxt.password, "mysql -u root -p", password, create, grant, quit)
}
