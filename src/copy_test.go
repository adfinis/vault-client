package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestCopy(t *testing.T) {

	var err error
	cfg, vc, err = SetupTestEnvironment()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	ui := cli.NewMockUi()
	c := &CopyCommand{Ui: ui}

	t.Run("TooFewArgs", func(t *testing.T) {

		args := []string{
			TestBackend + "/insertedsecret",
		}

		if rc := c.Run(args); rc != 1 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := "The copy command expects a source and a destination path"
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	t.Run("CopyNonexistentSourceSecret", func(t *testing.T) {

		args := []string{
			TestBackend + "/nonexistensecret",
			TestBackend + "/destinationsecret",
		}

		if rc := c.Run(args); rc != 1 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := "Source secret doesn't exist"
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	t.Run("CopyExistentSourceSecret", func(t *testing.T) {

		// Create test secret
		data := make(map[string]interface{})
		data["key"] = "value"

		_, err = vc.Logical().Write(TestBackend+"/existent", data)
		if err != nil {
			t.Fatalf("Unable to write test secret: %q", err)
		}

		args := []string{
			TestBackend + "/existent",
			TestBackend + "/destinationsecret",
		}

		if rc := c.Run(args); rc != 0 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := ""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	err = TeardownTestEnvironment()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
