package ssh

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/neofight/gostack"
	"golang.org/x/crypto/ssh"
)

func CreateSSHClient(host string, password []byte) (*ssh.Client, error) {

	currentUser, err := user.Current()

	if err != nil {
		return nil, fmt.Errorf("Unable to identify user: %v", err)
	}

	encryptedKey, err := ioutil.ReadFile(currentUser.HomeDir + "/.ssh/id_rsa")

	if err != nil {
		return nil, fmt.Errorf("Unable to read private key: %v", err)
	}

	var block, _ = pem.Decode(encryptedKey)

	decryptedKey, err := x509.DecryptPEMBlock(block, password)

	if err != nil {
		return nil, fmt.Errorf("Unable to decrypt private key: %v", err)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(decryptedKey)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse private key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)

	config := ssh.ClientConfig{
		User: currentUser.Username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}

	client, err := ssh.Dial("tcp", host+":22", &config)

	if err != nil {
		return nil, fmt.Errorf("Unable to connect: %v", err)
	}

	return client, nil
}

func ScpDownloadFile(client *ssh.Client, path string) (string, error) {

	// Ref: https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works

	session, err := client.NewSession()

	if err != nil {
		return "", fmt.Errorf("Unable to create session: %v", err)
	}

	defer session.Close()

	c := make(chan string)

	go func() {
		stdin, err := session.StdinPipe()

		if err != nil {
			return
		}

		stdout, err := session.StdoutPipe()

		if err != nil {
			return
		}

		readBuffer := make([]byte, 1024)

		fmt.Fprint(stdin, "\x00")

		n, err := stdout.Read(readBuffer)

		if err != nil && err != io.EOF {
			return
		}

		header := string(readBuffer[:n])

		size, err := strconv.Atoi(strings.Split(header, " ")[1])

		if err != nil {
			return
		}

		fmt.Fprint(stdin, "\x00")

		var read = 0

		var writeBuffer bytes.Buffer

		for {
			n, err := stdout.Read(readBuffer)

			if err != nil {
				if err == io.EOF {
					break
				} else {
					close(c)
					return
				}
			}

			if size-read < n {
				n = size - read
			}

			read += n

			writeBuffer.Write(readBuffer[:n])

			fmt.Fprint(stdin, "\x00")
		}

		c <- writeBuffer.String()
	}()

	err = session.Run("scp -f " + path)

	if err != nil {
		return "", err
	}

	result := <-c

	return result, nil
}

func ScpUploadData(client *ssh.Client, data string, filePath string, password []byte) error {

	// Ref: https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works

	_, filename := filepath.Split(filePath)

	inputs := make([]string, 3)

	inputs[0] = fmt.Sprintln("C0644", len(data), filename)
	inputs[1] = fmt.Sprint(data)
	inputs[2] = fmt.Sprint("\x00")

	_, err := runSudo(client, "scp -t "+filePath, inputs, password)

	if err != nil {
		return fmt.Errorf("Failed to upload to %v: %v", filePath, err)
	}

	return nil
}

type FilterFunc func(path string, info os.FileInfo) bool

func ScpUpload(client *ssh.Client, localPath string, remotePath string, password []byte, filter FilterFunc) error {

	// Ref: https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works

	inputs := make([]string, 0)

	var dirs gostack.Stack

	filepath.Walk(localPath, func(path string, info os.FileInfo, _ error) error {

		if !filter(path, info) {

			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		dir, _ := dirs.Peep()

		for dirs.Count() > 0 && !strings.HasPrefix(path, dir.(string)) {
			dirs.Pop()
			inputs = append(inputs, "E\n")

			dir, _ = dirs.Peep()
		}

		if info.IsDir() {
			dirs.Push(path)
			inputs = append(inputs, "D0755 0 "+info.Name()+"\n")
		} else {
			messages, err := createFileMessages(path)

			if err != nil {
				return fmt.Errorf("Unable to upload file: %v", err)
			}

			inputs = append(inputs, messages...)
		}

		return nil
	})

	_, err := runSudo(client, "scp -tr "+remotePath, inputs, password)

	if err != nil {
		return fmt.Errorf("Failed to complete upload: %v", err)
	}

	return nil
}

func createFileMessages(path string) ([]string, error) {

	_, filename := filepath.Split(path)

	buffer, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("Unable to read file: %v", err)
	}

	data := string(buffer)

	messages := make([]string, 3)

	messages[0] = fmt.Sprintln("C0644", len(data), filename)
	messages[1] = fmt.Sprint(data)
	messages[2] = fmt.Sprint("\x00")

	return messages, nil
}

func RunSudoCommands(client *ssh.Client, password []byte, commands ...string) (string, error) {

	var command string

	if len(commands) == 1 {
		command = commands[0]
	} else {
		command = fmt.Sprintf("sh -c '%v'", strings.Join(commands, ";"))
	}

	return runSudo(client, command, []string{}, password)
}

func RunCommand(client *ssh.Client, command string, inputs ...string) (string, error) {

	return run(client, command, inputs)
}

func runSudo(client *ssh.Client, command string, inputs []string, password []byte) (string, error) {

	command = "sudo -S " + command

	inputs = append([]string{string(password) + "\n"}, inputs...)

	return run(client, command, inputs)
}

func run(client *ssh.Client, command string, inputs []string) (string, error) {

	session, err := client.NewSession()

	if err != nil {
		return "", fmt.Errorf("Unable to create session: %v", err)
	}

	defer session.Close()

	var stdout bytes.Buffer
	session.Stdout = &stdout

	go func() {
		stdin, err := session.StdinPipe()

		if err != nil {
			return
		}

		defer stdin.Close()

		for _, input := range inputs {
			fmt.Fprint(stdin, input)
		}
	}()

	err = session.Run(command)

	if err != nil {
		return "", fmt.Errorf("There was an error executing the command %v: %v", command, err)
	}

	return stdout.String(), nil
}
