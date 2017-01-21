package main

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

func scpDownload(client *ssh.Client, path string) error {

	// Ref: https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works

	_, filename := filepath.Split(path)

	session, err := client.NewSession()

	if err != nil {
		return fmt.Errorf("Unable to create session: %v", err)
	}

	defer session.Close()

	go func() {
		stdin, err := session.StdinPipe()

		if err != nil {
			return
		}

		stdout, err := session.StdoutPipe()

		if err != nil {
			return
		}

		buffer := make([]byte, 1024)

		fmt.Fprint(stdin, "\x00")

		n, err := stdout.Read(buffer)

		if err != nil && err != io.EOF {
			return
		}

		header := string(buffer[:n])

		size, err := strconv.Atoi(strings.Split(header, " ")[1])

		if err != nil {
			return
		}

		fmt.Fprint(stdin, "\x00")

		file, err := os.Create(filename)

		if err != nil {
			return
		}

		defer file.Close()

		var read = 0

		for {
			n, err := stdout.Read(buffer)

			if err != nil {
				if err == io.EOF {
					break
				} else {
					return
				}
			}

			if size-read < n {
				n = size - read
			}

			read += n

			file.Write(buffer[:n])

			fmt.Fprint(stdin, "\x00")
		}
	}()

	err = session.Run("scp -f " + path)

	return err
}

func scpUpload(client *ssh.Client, data string, filename string) error {

	// Ref: https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works

	session, err := client.NewSession()

	if err != nil {
		return fmt.Errorf("Unable to create session: %v", err)
	}

	defer session.Close()

	go func() {
		stdin, err := session.StdinPipe()

		if err != nil {
			return
		}

		defer stdin.Close()

		fmt.Fprintln(stdin, "C0644", len(data), filename)
		fmt.Fprint(stdin, data)
		fmt.Fprint(stdin, "\x00")
	}()

	_, err = runCommands(session, "scp -t "+filename)

	return err
}

func runCommands(session *ssh.Session, commands ...string) (string, error) {

	var buffer bytes.Buffer
	session.Stdout = &buffer

	command := strings.Join(commands, ";")

	err := session.Run(command)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func runSudoCommands(session *ssh.Session, password []byte, commands ...string) (string, error) {

	go func() {
		stdin, err := session.StdinPipe()

		if err != nil {
			return
		}

		defer stdin.Close()

		fmt.Fprintln(stdin, string(password))
	}()

	var buffer bytes.Buffer
	session.Stdout = &buffer

	command := fmt.Sprintf("sudo -S sh -c '%v'", strings.Join(commands, ";"))

	err := session.Run(command)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
