package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var stdOut io.Writer = os.Stdout

const usage string = "Usage: govps <host> <domain>"

type arguments struct {
	host   string
	domain string
}

func parse(args []string) (parsed arguments, ok bool) {

	if len(args) != 3 {
		printUsage()
		return arguments{}, false
	}

	host := strings.TrimSpace(args[1])
	domain := strings.TrimSpace(args[2])

	if !isArgumentValid(host) || !isArgumentValid(domain) {
		printUsage()
		return arguments{}, false
	}

	return arguments{
		host:   host,
		domain: domain,
	}, true
}

func isArgumentValid(argument string) bool {

	match, _ := regexp.MatchString("^[A-Za-z0-9.\\-]+$", argument)

	return match
}

func printUsage() {
	fmt.Fprintln(stdOut, usage)
}
