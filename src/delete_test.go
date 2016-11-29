package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestDelete_TooManyArgs(t *testing.T) {

	ui := new(cli.MockUi)
	c := &DeleteCommand{
		Ui: ui,
	}

	args := []string{
		"secret/doesntexist",
		"secret/toomucharguments",
	}

	if rc := c.Run(args); rc != 1 {
		t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
	}

	expected := "The rm command expects at most one argument"
	if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
		t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
	}

}

func TestDelete_NonexistentSecret(t *testing.T) {

	err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	err = InitializeClient(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	ui := new(cli.MockUi)
	c := &DeleteCommand{Ui: ui}

	args := []string{"secret/doesntexist"}

	if rc := c.Run(args); rc != 1 {
		t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
	}

	expected := "Secret does not exist"
	if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
		t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
	}

}
