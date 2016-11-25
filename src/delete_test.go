package main

import (
	"strings"
	"testing"

	"fmt"
	"github.com/mitchellh/cli"
)

func TestDelete_tooManyArgs(t *testing.T) {

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
