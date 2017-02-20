package tasks

import (
	"fmt"
	"testing"

	"github.com/neofight/govps/io"
	"github.com/neofight/govps/mock"
)

const webConfigContents = `<?xml version="1.0"?>
<configuration>
  <connectionStrings>
    <add name="mysql" connectionString="Server=localhost;Database=TestDB;Uid=TestUser;Pwd=TestPassword;" providerName="MySql.Data.MySqlClient" />
  </connectionStrings>
</configuration>`

func sqlWasRun(server *mock.Server, sql string) bool {

	for _, command := range server.CommandsRun {

		if command.Command != mysqlCommand {
			continue
		}

		if contains(command.Inputs, sql) {
			return true
		}
	}

	return false
}

func TestCreateMySQLDatabaseNoDatabase(t *testing.T) {

	mockFS := mock.NewFileSystem()

	mockFS.Files["Web.config"] = webConfigContents

	io.FileSystem = mockFS

	server := mock.NewServer()

	cxt := Context{server, "test.com"}

	mysql := CreateMySQLDatabase{[]byte("password")}

	err := mysql.Execute(cxt)

	if err != nil {
		t.Error("Expected task to run without error but it did not")
	}

	if !sqlWasRun(server, fmt.Sprintf(createDatabaseCommand, "TestDB")) {
		t.Error("Expected database to be created but it was not")
	}
}

func TestCreateMySQLDatabaseExistingDatabase(t *testing.T) {

	mockFS := mock.NewFileSystem()

	mockFS.Files["Web.config"] = webConfigContents

	io.FileSystem = mockFS

	server := mock.NewServer()

	server.Responses[fmt.Sprintf(checkDatabaseExistsCommand, "TestDB")] = "TestDB"

	cxt := Context{server, "test.com"}

	mysql := CreateMySQLDatabase{[]byte("password")}

	err := mysql.Execute(cxt)

	if err != nil {
		t.Error("Expected task to run without error but it did not")
	}

	if sqlWasRun(server, fmt.Sprintf(createDatabaseCommand, "TestDB")) {
		t.Error("Expected database not to be created but it was")
	}
}
