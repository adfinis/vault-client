package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestList(t *testing.T) {

	err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	err = InitializeClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	ui := cli.NewMockUi()
	c := &ListCommand{Ui: ui}

	t.Run("ListSecretsInEmptyBackend", func(t *testing.T) {

		args := []string{"secret/"}

		if rc := c.Run(args); rc != 0 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expectedErr := ""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expectedErr) {
			t.Fatalf("expected error:\n%s\n\nto include: %q", actual, expectedErr)
		}

		expectedOutput := ""
		if actual := ui.OutputWriter.String(); !strings.Contains(actual, expectedOutput) {
			t.Fatalf("expected output:\n%s\n\nto include: %q", actual, expectedOutput)
		}
	})

	t.Run("ListExistingSecrets", func(t *testing.T) {

		data := make(map[string]interface{})
		data["key"] = "value"

		for i := 1; i <= 3; i++ {
			_, err = kv.Put(fmt.Sprintf("secret/secret%v", i), data)
			if err != nil {
				t.Fatalf("Unable to write test secret: %q", err)
			}
		}

		args := []string{"secret/"}

		if rc := c.Run(args); rc != 0 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expectedErr := ""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expectedErr) {
			t.Fatalf("expected error:\n%s\n\nto include: %q", actual, expectedErr)
		}

		expectedOutput := `secret1
secret2
secret3`
		if actual := ui.OutputWriter.String(); !strings.Contains(actual, expectedOutput) {
			t.Fatalf("expected output:\n%s\n\nto include: %q", actual, expectedOutput)
		}

	})

	t.Run("ListExistingSecretsRecusively", func(t *testing.T) {

		args := []string{"-r", ""}

		if rc := c.Run(args); rc != 0 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expectedErr := ""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expectedErr) {
			t.Fatalf("expected error:\n%s\n\nto include: %q", actual, expectedErr)
		}

		expectedOutput := `secret/secret1
secret/secret2
secret/secret3`
		if actual := ui.OutputWriter.String(); !strings.Contains(actual, expectedOutput) {
			t.Fatalf("expected output:\n%s\n\nto include: %q", actual, expectedOutput)
		}

	})

	for i := 1; i <= 3; i++ {
		_, err = kv.Delete(fmt.Sprintf("secret/secret%v", i))
		if err != nil {
			t.Fatalf("Unable to write test secret: %q", err)
		}
	}

	_, err = kv.Delete("secret/directory/secret1")
	if err != nil {
		t.Fatalf("Unable to write test secret: %q", err)
	}
}
