package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestInsert(t *testing.T) {

	var err error
	cfg, vc, err = SetupTestEnvironment()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	ui := cli.NewMockUi()
	c := &InsertCommand{Ui: ui}

	t.Run("TooFewArgs", func(t *testing.T) {

		args := []string{
			TestBackend + "/insertedsecret",
		}

		if rc := c.Run(args); rc != 1 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := "The insert command expects at least a path and a k/v pair (key=value)"
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	t.Run("InsertInvalidKvArgs", func(t *testing.T) {

		args := []string{
			TestBackend + "/insertedsecret",
			"invalidkey: invalidvalue",
		}

		if rc := c.Run(args); rc != 1 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := "Invalid key/value arguments: \"invalidkey: invalidvalue\""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	t.Run("InsertValidKvArgs", func(t *testing.T) {

		args := []string{
			TestBackend + "/insertedsecret",
			"invalidkey=invalidvalue",
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
