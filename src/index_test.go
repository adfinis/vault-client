package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestIndex(t *testing.T) {

	err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	err = InitializeClient(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	ui := new(cli.MockUi)
	c := &IndexCommand{Ui: ui}

	t.Run("IndexEmptyVault", func(t *testing.T) {

		if rc := c.Run([]string{}); rc != 1 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := "Unable to index vault: \"Backend \\\"secret/\\\" holds no secrets\""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	t.Run("IndexVault", func(t *testing.T) {

		// Create test secret
		data := make(map[string]interface{})
		data["key"] = "value"

		_, err = vc.Logical().Write("secret/indexsecret", data)
		if err != nil {
			t.Fatalf("Unable to write test secret: %q", err)
		}

		if rc := c.Run([]string{}); rc != 0 {
			t.Fatalf("Wrong exit code. errors: \n%s", ui.ErrorWriter.String())
		}

		expected := ""
		if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
			t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
		}
	})

	_, err = vc.Logical().Delete("secret/indexsecret")
	if err != nil {
		t.Fatalf("Unable to clean up test secret: %q", err)
	}
}
