package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/neofight/govps/io"
)

const Usage string = "Usage: govps <host> <domain>"

type Arguments struct {
	Host   string
	Domain string
}

func Parse(args []string) (parsed Arguments, ok bool) {

	if len(args) != 3 {
		printUsage()
		return Arguments{}, false
	}

	host := strings.TrimSpace(args[1])
	domain := strings.TrimSpace(args[2])

	if !isArgumentValid(host) || !isArgumentValid(domain) {
		printUsage()
		return Arguments{}, false
	}

	return Arguments{
		Host:   host,
		Domain: domain,
	}, true
}

func isArgumentValid(argument string) bool {

	match, _ := regexp.MatchString("^[A-Za-z0-9.\\-]+$", argument)

	return match
}

func printUsage() {
	fmt.Fprintln(io.StdOut, Usage)
}
