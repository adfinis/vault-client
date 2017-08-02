package main

import (
	"fmt"
	"os"
	"testing"
)

func TestAuth(t *testing.T) {

	err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	t.Run("VaultUnreachable", func(t *testing.T) {

		cfg.Host = "127.0.0.1"
		cfg.Port = 33333

		err = InitializeClient()
		if err == nil {
			t.Fatalf("Test expects error to occur")
		}

		expected := "Unable to connect to vault"
		if actual := err.Error(); actual != expected {
			t.Fatalf("expected:\n%s\n\nto be: %q", actual, expected)
		}
	})
}
