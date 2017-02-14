package main

import (
	"fmt"

	"github.com/neofight/govps/io"
)

func promptForPassword(which string) ([]byte, error) {

	fmt.Printf("Enter %v password:", which)
	password, err := io.Terminal.ReadPassword(0)
	fmt.Println()

	return password, err
}
