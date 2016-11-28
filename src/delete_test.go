package main

import (
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

/*func TestDelete_ExistentSecret(t *testing.T) {

	ui := new(cli.MockUi)
	c := &DeleteCommand{Ui: ui}

	args := []string{"secret/exists"}

	// Create secret so it exists
	data := make(map[string]interface{})
	_, err := vc.Logical().Write(args[0], data)
	if err != nil {
		t.Fatalf("Unable to write example feature: %q", err)
	}

	if rc := c.Run(args); rc != 1 {
		t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
	}

	expected := "Secret does not exist"
	if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
		t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
	}

}*/
