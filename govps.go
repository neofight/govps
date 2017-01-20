package main

import "bytes"
import "crypto/x509"
import "encoding/pem"
import "fmt"
import "golang.org/x/crypto/ssh"
import "golang.org/x/crypto/ssh/terminal"
import "io/ioutil"
import "log"
import "os"
import "os/user"

func main() {
	host := os.Args[1]

	currentUser, err := user.Current()

	if err != nil {
		log.Fatal("Unable to identify user: ", err)
	}

	encryptedKey, err := ioutil.ReadFile(currentUser.HomeDir + "/.ssh/id_rsa")

	if err != nil {
		log.Fatal("Unable to read private key: ", err)
	}

	fmt.Print("Enter password:")
	password, err := terminal.ReadPassword(0)
	fmt.Println()

	if err != nil {
		log.Fatal("Unable to read password: ", err)
	}

	var block, _ = pem.Decode(encryptedKey)

	decryptedKey, err := x509.DecryptPEMBlock(block, password)

	if err != nil {
		log.Fatal("Unable to decrypt private key: ", err)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(decryptedKey)

	if err != nil {
		log.Fatal("Unable to parse private key: ", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)

	config := ssh.ClientConfig{
		User: currentUser.Username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}

	client, err := ssh.Dial("tcp", host+":22", &config)

	if err != nil {
		log.Fatal("Unable to connect: ", err)
	}

	defer client.Close()

	session, err := client.NewSession()

	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}

	defer session.Close()

	var buffer bytes.Buffer
	session.Stdout = &buffer
	err = session.Run("uname -a")

	if err != nil {
		log.Fatal("Failed to run remote command: ", err)
	}

	fmt.Print(buffer.String())
}
