package main

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

func promptForPassword(which string) ([]byte, error) {

	fmt.Printf("Enter %v password:", which)
	password, err := terminal.ReadPassword(0)
	fmt.Println()

	return password, err
}
