package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestMove(t *testing.T) {

	err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	err = InitializeClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	ui := new(cli.MockUi)
	c := &MoveCommand{Ui: ui}

	t.Run("TooFewArgs", func(t *testing.T) {

		args := []string{
			"secret/insertedsecret",
		}

		if rc := c.Run(args); rc != 1 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := "The move command expects a source and a destination path"
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	t.Run("MoveNonexistentSourceSecret", func(t *testing.T) {

		args := []string{
			"secret/nonexistensecret",
			"secret/destinationsecret",
		}

		if rc := c.Run(args); rc != 1 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := "Source secret doesn't exist"
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	t.Run("MoveExistentSourceSecret", func(t *testing.T) {

		// Create test secret
		data := make(map[string]interface{})
		data["key"] = "value"

		_, err = vc.Logical().Write("secret/existent", data)
		if err != nil {
			t.Fatalf("Unable to write test secret: %q", err)
		}

		args := []string{
			"secret/existent",
			"secret/destinationsecret",
		}

		if rc := c.Run(args); rc != 0 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := ""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	_, err = vc.Logical().Delete("secret/destinationsecret")
	if err != nil {
		t.Fatalf("Unable to clean up test secret: %q", err)
	}
}
