package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func createSSHClient(host string, password []byte) (*ssh.Client, error) {

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

func scpDownloadFile(client *ssh.Client, path string) (string, error) {

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

func scpUploadDataAsUser(client *ssh.Client, data string, filePath string) error {

	return scpUploadData(client, data, filePath, run)
}

func scpUploadDataAsRoot(client *ssh.Client, data string, filePath string, password []byte) error {

	return scpUploadData(client, data, filePath, func(session *ssh.Session, command string, inputs []string) error {

		return runSudo(session, command, inputs, password)
	})
}

func scpUploadData(client *ssh.Client, data string, filePath string, run func(*ssh.Session, string, []string) error) error {

	// Ref: https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works

	_, filename := filepath.Split(filePath)

	session, err := client.NewSession()

	if err != nil {
		return fmt.Errorf("Unable to create session: %v", err)
	}

	defer session.Close()

	inputs := make([]string, 3)

	inputs[0] = fmt.Sprintln("C0644", len(data), filename)
	inputs[1] = fmt.Sprint(data)
	inputs[2] = fmt.Sprint("\x00")

	return run(session, "scp -t "+filePath, inputs)
}

func runSudoCommand(session *ssh.Session, command string, password []byte) error {

	return runSudo(session, command, []string{}, password)
}

func runSudo(session *ssh.Session, command string, inputs []string, password []byte) error {

	command = "sudo -S " + command

	inputs = append([]string{string(password) + "\n"}, inputs...)

	return run(session, command, inputs)
}

func run(session *ssh.Session, command string, inputs []string) error {

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

	return session.Run(command)
}
