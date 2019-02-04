package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestShow(t *testing.T) {

	var err error
	cfg, vc, err = SetupTestEnvironment()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	ui := cli.NewMockUi()
	c := &ShowCommand{Ui: ui}

	t.Run("ShowNonexistentSecret", func(t *testing.T) {

		args := []string{TestBackend + "/secret2"}

		if rc := c.Run(args); rc != 1 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expectedErr := "Secret does not exist"
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expectedErr) {
			t.Fatalf("expected error:\n%s\n\nto include: %q", actual, expectedErr)
		}

		expectedOutput := ""
		if actual := ui.OutputWriter.String(); !strings.Contains(actual, expectedOutput) {
			t.Fatalf("expected output:\n%s\n\nto include: %q", actual, expectedOutput)
		}
	})

	t.Run("ShowExistentSecret", func(t *testing.T) {

		// Create test secret
		data := make(map[string]interface{})
		data["key"] = "value"

		_, err = vc.Logical().Write(TestBackend+"/secret1", data)
		if err != nil {
			t.Fatalf("Unable to write test secret: %q", err)
		}

		args := []string{TestBackend + "/secret1"}

		if rc := c.Run(args); rc != 0 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expectedErr := ""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expectedErr) {
			t.Fatalf("expected error:\n%s\n\nto include: %q", actual, expectedErr)
		}

		expectedOutput := "key: value"
		if actual := ui.OutputWriter.String(); !strings.Contains(actual, expectedOutput) {
			t.Fatalf("expected output:\n%s\n\nto include: %q", actual, expectedOutput)
		}
	})

	t.Run("ShowSortedSecrets", func(t *testing.T) {

		// Create test secret
		data := make(map[string]interface{})
		data["a_key"] = "value"
		data["c_key"] = "value"
		data["b_key"] = "value"

		_, err = vc.Logical().Write(TestBackend+"/secret1", data)
		if err != nil {
			t.Fatalf("Unable to write test secret: %q", err)
		}

		args := []string{TestBackend + "/secret1"}

		if rc := c.Run(args); rc != 0 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expectedErr := ""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expectedErr) {
			t.Fatalf("expected error:\n%s\n\nto include: %q", actual, expectedErr)
		}

		expectedOutput := `a_key: value
b_key: value
c_key: value`
		if actual := ui.OutputWriter.String(); !strings.Contains(actual, expectedOutput) {
			t.Fatalf("expected output:\n%s\n\nto include: %q", actual, expectedOutput)
		}
	})

	err = TeardownTestEnvironment()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
