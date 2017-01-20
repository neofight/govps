package main

import "bytes"
import "crypto/x509"
import "encoding/pem"
import "fmt"
import "golang.org/x/crypto/ssh"
import "golang.org/x/crypto/ssh/terminal"
import "io/ioutil"
import "os/user"
import "strings"

type sshClient struct {
	*ssh.Client
	*ssh.Session
}

func createSSHClient(host string) (*sshClient, error) {

	currentUser, err := user.Current()

	if err != nil {
		return nil, fmt.Errorf("Unable to identify user: %v", err)
	}

	encryptedKey, err := ioutil.ReadFile(currentUser.HomeDir + "/.ssh/id_rsa")

	if err != nil {
		return nil, fmt.Errorf("Unable to read private key: %v", err)
	}

	fmt.Print("Enter password:")
	password, err := terminal.ReadPassword(0)
	fmt.Println()

	if err != nil {
		return nil, fmt.Errorf("Unable to read password: %v", err)
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

	session, err := client.NewSession()

	if err != nil {
		client.Close()
		return nil, fmt.Errorf("Failed to create session: %v", err)
	}

	return &sshClient{client, session}, nil
}

func (client sshClient) close() {
	client.Client.Close()
	client.Session.Close()
}

func (client sshClient) runCommands(commands ...string) (string, error) {
	var buffer bytes.Buffer
	client.Session.Stdout = &buffer

	command := strings.Join(commands, ";")

	err := client.Session.Run(command)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
