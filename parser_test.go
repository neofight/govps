package main

import (
	"bytes"
	"testing"

	"github.com/neofight/govps/io"
)

func TestParserInsufficientArguments(t *testing.T) {

	buffer := &bytes.Buffer{}
	io.StdOut = buffer

	_, ok := parse([]string{"govps", "one"})

	if ok {
		t.Error("Expected command with one argument to be rejected but was accepted")
	}

	if buffer.String() != usage+"\n" {
		t.Error("Expected command help to be printed but it was not")
	}
}

func TestParserExcessiveArguments(t *testing.T) {

	buffer := &bytes.Buffer{}
	io.StdOut = buffer

	_, ok := parse([]string{"govps", "one", "two", "three"})

	if ok {
		t.Error("Expected command with three arguments to be rejected but was accepted")
	}

	if buffer.String() != usage+"\n" {
		t.Error("Expected command help to be printed but it was not")
	}

	if buffer.String() != usage+"\n" {
		t.Error("Expected command help to be printed but it was not")
	}
}

var illegalCharacters = []string{";", "'", "\"", "\r", "\n", "`"}

func TestParserIllegalCharactersInHost(t *testing.T) {

	buffer := &bytes.Buffer{}
	io.StdOut = buffer

	for _, character := range illegalCharacters {

		_, ok := parse([]string{"govps", character, "test.com"})

		if ok {
			t.Errorf("Expected host %v to be rejected but was accepted", character)

			if buffer.String() != usage+"\n" {
				t.Error("Expected command help to be printed but it was not")
			}

			buffer.Reset()
		}
	}
}

func TestParserIllegalCharactersInDomain(t *testing.T) {

	buffer := &bytes.Buffer{}
	io.StdOut = buffer

	for _, character := range illegalCharacters {

		_, ok := parse([]string{"govps", "test.com", character})

		if ok {
			t.Errorf("Expected domain %v to be rejected but was accepted", character)
		}

		if buffer.String() != usage+"\n" {
			t.Error("Expected command help to be printed but it was not")
		}

		buffer.Reset()
	}
}

func TestParserHappyPath(t *testing.T) {

	buffer := &bytes.Buffer{}
	io.StdOut = buffer

	result, ok := parse([]string{"govps", "host.com", "domain.com"})

	if !ok {
		t.Fatal("Expected commmand to be parsed but it was not")
	}

	if result.host != "host.com" {
		t.Errorf("Expected host host.com but saw %v", result.host)
	}

	if result.domain != "domain.com" {
		t.Errorf("Expected domain domain.com but saw %v", result.domain)
	}

	if buffer.String() != "" {
		t.Errorf("Expected no output but some was printed:\n%v", buffer.String())
	}
}
