package main

import (
	"fmt"
)

func promptForPassword(which string) ([]byte, error) {

	fmt.Printf("Enter %v password:", which)
	password, err := term.ReadPassword(0)
	fmt.Println()

	return password, err
}
